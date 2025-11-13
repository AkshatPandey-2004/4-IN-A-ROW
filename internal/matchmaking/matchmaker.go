package matchmaking

import (
    "sync"
    "time"
    "github.com/AkshatPandey-2004/4-IN-A-ROW/pkg/models"
    "github.com/AkshatPandey-2004/4-IN-A-ROW/internal/game"
)

type WaitingPlayer struct {
    Player    *models.Player
    Timestamp time.Time
    GameChan  chan *game.GameInstance
}

type Matchmaker struct {
    waiting map[string]*WaitingPlayer
    mu      sync.Mutex
}

func NewMatchmaker() *Matchmaker {
    return &Matchmaker{
        waiting: make(map[string]*WaitingPlayer),
    }
}

func (m *Matchmaker) AddPlayer(player *models.Player) chan *game.GameInstance {
    m.mu.Lock()
    defer m.mu.Unlock()

    gameChan := make(chan *game.GameInstance, 1)
    
    m.waiting[player.ID] = &WaitingPlayer{
        Player:    player,
        Timestamp: time.Now(),
        GameChan:  gameChan,
    }

    return gameChan
}

func (m *Matchmaker) RemovePlayer(playerID string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.waiting, playerID)
}

func (m *Matchmaker) TryMatch(playerID string) (*game.GameInstance, *models.Player) {
    m.mu.Lock()
    defer m.mu.Unlock()

    wp, exists := m.waiting[playerID]
    if !exists {
        return nil, nil
    }

    // Try to find another waiting player
    for id, otherWP := range m.waiting {
        if id != playerID {
            // Found a match!
            player1 := wp.Player
            player2 := otherWP.Player
            
            player1.Piece = 1
            player2.Piece = 2

            newGame := game.NewGame(player1, false)
            newGame.AddPlayer2(player2)

            // Notify both players
            wp.GameChan <- newGame
            otherWP.GameChan <- newGame

            // Remove both from waiting
            delete(m.waiting, playerID)
            delete(m.waiting, id)

            return newGame, player2
        }
    }

    return nil, nil
}

func (m *Matchmaker) GetWaitingPlayer(playerID string) *WaitingPlayer {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.waiting[playerID]
}