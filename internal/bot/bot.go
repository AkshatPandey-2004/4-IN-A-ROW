package bot

import (
    "math/rand"
    "github.com/AkshatPandey-2004/4-IN-A-ROW/internal/game"
)

type Bot struct {
    playerNum int
}

func NewBot(playerNum int) *Bot {
    // Removed rand.Seed() - Go 1.20+ auto-seeds the random generator
    return &Bot{playerNum: playerNum}
}

func (b *Bot) GetMove(board *game.Board) int {
    // Strategy priority:
    // 1. Win if possible
    // 2. Block opponent from winning
    // 3. Strategic move (center preference)
    
    availableCols := board.GetAvailableColumns()
    if len(availableCols) == 0 {
        return -1
    }

    opponent := 1
    if b.playerNum == 1 {
        opponent = 2
    }

    // Check if bot can win
    for _, col := range availableCols {
        if b.isWinningMove(board, col, b.playerNum) {
            return col
        }
    }

    // Block opponent from winning
    for _, col := range availableCols {
        if b.isWinningMove(board, col, opponent) {
            return col
        }
    }

    // Strategic move - prefer center columns
    centerCols := []int{3, 2, 4, 1, 5, 0, 6}
    for _, col := range centerCols {
        for _, availCol := range availableCols {
            if col == availCol {
                return col
            }
        }
    }

    // Fallback to random
    return availableCols[rand.Intn(len(availableCols))]
}

func (b *Bot) isWinningMove(board *game.Board, col int, player int) bool {
    clonedBoard := board.Clone()
    row, ok := clonedBoard.MakeMove(col, player)
    if !ok {
        return false
    }
    return clonedBoard.CheckWin(row, col, player)
}