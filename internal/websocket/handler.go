package websocket

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
    "github.com/gorilla/websocket"
    "github.com/google/uuid"
    "github.com/AkshatPandey-2004/4-in-a-row/pkg/models"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/game"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/bot"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/matchmaking"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/database"
    "github.com/AkshatPandey-2004/4-in-a-row/internal/kafka"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Handler struct {
    hub         *Hub
    gameManager *game.Manager
    matchmaker  *matchmaking.Matchmaker
    db          *database.DB
    kafkaProducer *kafka.Producer
}

func NewHandler(hub *Hub, gameManager *game.Manager, matchmaker *matchmaking.Matchmaker, 
    db *database.DB, kafkaProducer *kafka.Producer) *Handler {
    return &Handler{
        hub:         hub,
        gameManager: gameManager,
        matchmaker:  matchmaker,
        db:          db,
        kafkaProducer: kafkaProducer,
    }
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    username := r.URL.Query().Get("username")
    if username == "" {
        conn.Close()
        return
    }

    playerID := uuid.New().String()
    client := NewClient(playerID, username, h.hub, conn)
    
    h.hub.register <- client

    // Save player to database
    h.db.CreateOrGetPlayer(username, playerID)

    go client.writePump()
    go client.readPump(h.handleMessage)
}

func (h *Handler) handleMessage(client *Client, message []byte) {
    var msg map[string]interface{}
    if err := json.Unmarshal(message, &msg); err != nil {
        log.Println("Error parsing message:", err)
        return
    }

    msgType, ok := msg["type"].(string)
    if !ok {
        return
    }

    switch msgType {
    case "find_match":
        h.handleFindMatch(client)
    case "make_move":
        h.handleMakeMove(client, msg)
    case "rejoin":
        h.handleRejoin(client, msg)
    }
}

func (h *Handler) handleFindMatch(client *Client) {
    player := &models.Player{
        ID:       client.id,
        Username: client.username,
        Piece:    1,
    }

    gameChan := h.matchmaker.AddPlayer(player)

    // Try immediate match
    gameInstance, opponent := h.matchmaker.TryMatch(player.ID)
    
    if gameInstance != nil {
        // Found a match
        h.startGame(client, gameInstance, opponent)
        return
    }

    // Wait for match with timeout
    go func() {
        select {
        case game := <-gameChan:
            h.startGame(client, game, nil)
        case <-time.After(10 * time.Second):
            // Timeout - start game with bot
            h.startGameWithBot(client, player)
        }
    }()
}

func (h *Handler) startGame(client *Client, gameInstance *game.GameInstance, opponent *models.Player) {
    client.gameID = gameInstance.ID
    h.gameManager.AddGame(gameInstance)

    // Send game start to client
    client.SendJSON(map[string]interface{}{
        "type": "game_start",
        "game": gameInstance.Game,
    })

    // Send to opponent if exists
    if opponent != nil {
        oppClient := h.hub.GetClient(opponent.ID)
        if oppClient != nil {
            oppClient.gameID = gameInstance.ID
            oppClient.SendJSON(map[string]interface{}{
                "type": "game_start",
                "game": gameInstance.Game,
            })
        }
    }

    // Send analytics event
    h.kafkaProducer.SendGameEvent(&models.GameEvent{
        Type:      "game_start",
        GameID:    gameInstance.ID,
        Data:      gameInstance.Game,
        Timestamp: time.Now(),
    })
}

func (h *Handler) startGameWithBot(client *Client, player *models.Player) {
    h.matchmaker.RemovePlayer(player.ID)

    botPlayer := &models.Player{
        ID:       "bot-" + uuid.New().String(),
        Username: "Bot",
        Piece:    2,
    }

    gameInstance := game.NewGame(player, true)
    gameInstance.AddPlayer2(botPlayer)
    
    client.gameID = gameInstance.ID
    h.gameManager.AddGame(gameInstance)

    client.SendJSON(map[string]interface{}{
        "type": "game_start",
        "game": gameInstance.Game,
    })

    h.kafkaProducer.SendGameEvent(&models.GameEvent{
        Type:      "game_start",
        GameID:    gameInstance.ID,
        Data:      gameInstance.Game,
        Timestamp: time.Now(),
    })
}

func (h *Handler) handleMakeMove(client *Client, msg map[string]interface{}) {
    gameID := client.gameID
    col, ok := msg["column"].(float64)
    if !ok {
        return
    }

    gameInstance, exists := h.gameManager.GetGame(gameID)
    if !exists {
        return
    }

    playerNum := 1
    if client.id == gameInstance.Player2.ID {
        playerNum = 2
    }

    row, success, result := gameInstance.MakeMove(int(col), playerNum)
    if !success {
        client.SendJSON(map[string]interface{}{
            "type":  "error",
            "message": result,
        })
        return
    }

    // Broadcast move to both players
    moveData := map[string]interface{}{
        "type":   "move_made",
        "column": int(col),
        "row":    row,
        "player": playerNum,
        "game":   gameInstance.Game,
    }

    client.SendJSON(moveData)
    
    // Send to opponent
    opponentID := gameInstance.Player2.ID
    if playerNum == 2 {
        opponentID = gameInstance.Player1.ID
    }
    
    if opponentClient := h.hub.GetClient(opponentID); opponentClient != nil {
        opponentClient.SendJSON(moveData)
    }

    // Send analytics event
    h.kafkaProducer.SendGameEvent(&models.GameEvent{
        Type:      "move_made",
        GameID:    gameID,
        Data:      moveData,
        Timestamp: time.Now(),
    })

    // Handle game end
    if result == "win" || result == "draw" {
        h.handleGameEnd(gameInstance, result)
        
        // Bot's turn if game continues and it's bot game
        if gameInstance.IsBot && result == "continue" && gameInstance.CurrentTurn == 2 {
            time.AfterFunc(500*time.Millisecond, func() {
                h.makeBotMove(gameInstance)
            })
        }
    }
}

func (h *Handler) makeBotMove(gameInstance *game.GameInstance) {
    botAI := bot.NewBot(2)
    col := botAI.GetMove(gameInstance.GetBoard())
    
    if col == -1 {
        return
    }

    row, success, result := gameInstance.MakeMove(col, 2)
    if !success {
        return
    }

    // Send bot move to player
    moveData := map[string]interface{}{
        "type":   "move_made",
        "column": col,
        "row":    row,
        "player": 2,
        "game":   gameInstance.Game,
    }

    if playerClient := h.hub.GetClient(gameInstance.Player1.ID); playerClient != nil {
        playerClient.SendJSON(moveData)
    }

    h.kafkaProducer.SendGameEvent(&models.GameEvent{
        Type:      "move_made",
        GameID:    gameInstance.ID,
        Data:      moveData,
        Timestamp: time.Now(),
    })

    if result == "win" || result == "draw" {
        h.handleGameEnd(gameInstance, result)
    }
}

func (h *Handler) handleGameEnd(gameInstance *game.GameInstance, result string) {
    // Save to database
    h.db.SaveGame(gameInstance.Game)
    h.db.UpdateGameStats(gameInstance.Game)

    // Send game end event
    endData := map[string]interface{}{
        "type":   "game_end",
        "result": result,
        "winner": gameInstance.Winner,
        "game":   gameInstance.Game,
    }

    if player1Client := h.hub.GetClient(gameInstance.Player1.ID); player1Client != nil {
        player1Client.SendJSON(endData)
    }

    if gameInstance.Player2 != nil {
        if player2Client := h.hub.GetClient(gameInstance.Player2.ID); player2Client != nil {
            player2Client.SendJSON(endData)
        }
    }

    h.kafkaProducer.SendGameEvent(&models.GameEvent{
        Type:      "game_end",
        GameID:    gameInstance.ID,
        Data:      endData,
        Timestamp: time.Now(),
    })
}

func (h *Handler) handleRejoin(client *Client, msg map[string]interface{}) {
    gameID, ok := msg["game_id"].(string)
    if !ok {
        return
    }

    gameInstance, exists := h.gameManager.GetGame(gameID)
    if !exists {
        client.SendJSON(map[string]interface{}{
            "type":    "error",
            "message": "Game not found",
        })
        return
    }

    client.gameID = gameID
    client.SendJSON(map[string]interface{}{
        "type": "rejoin_success",
        "game": gameInstance.Game,
    })
}