package game

import (
    "time"
    "github.com/AkshatPandey-2004/4-in-a-row/pkg/models"
    "github.com/google/uuid"
)

type GameInstance struct {
    *models.Game
    board *Board
}

func NewGame(player1 *models.Player, isBot bool) *GameInstance {
    board := NewBoard()
    game := &models.Game{
        ID:          uuid.New().String(),
        Player1:     player1,
        Board:       board.GetGrid(),
        CurrentTurn: 1,
        Status:      models.StatusWaiting,
        IsBot:       isBot,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    return &GameInstance{
        Game:  game,
        board: board,
    }
}

func (g *GameInstance) AddPlayer2(player2 *models.Player) {
    g.Player2 = player2
    g.Status = models.StatusPlaying
    g.UpdatedAt = time.Now()
}

func (g *GameInstance) MakeMove(col int, playerNum int) (int, bool, string) {
    if g.Status != models.StatusPlaying {
        return -1, false, "game is not in playing state"
    }

    if playerNum != g.CurrentTurn {
        return -1, false, "not your turn"
    }

    row, ok := g.board.MakeMove(col, playerNum)
    if !ok {
        return -1, false, "invalid move"
    }

    g.Board = g.board.GetGrid()
    g.UpdatedAt = time.Now()

    // Check for win
    if g.board.CheckWin(row, col, playerNum) {
        g.Status = models.StatusFinished
        if playerNum == 1 {
            g.Winner = g.Player1
        } else {
            g.Winner = g.Player2
        }
        now := time.Now()
        g.FinishedAt = &now
        return row, true, "win"
    }

    // Check for draw
    if g.board.IsFull() {
        g.Status = models.StatusFinished
        now := time.Now()
        g.FinishedAt = &now
        return row, true, "draw"
    }

    // Switch turn
    if g.CurrentTurn == 1 {
        g.CurrentTurn = 2
    } else {
        g.CurrentTurn = 1
    }

    return row, true, "continue"
}

func (g *GameInstance) GetBoard() *Board {
    return g.board
}