package bot

import (
	"math/rand"
	"time"

	"github.com/AkshatPandey-2004/4-IN-A-ROW/internal/game"
)

type Bot struct {
	playerNum int
}

func NewBot(playerNum int) *Bot {
	return &Bot{playerNum: playerNum}
}

// GetMove returns the best move for the bot with intentional mistakes for balance
func (b *Bot) GetMove(board *game.Board) int {
	// Add 2-second delay before bot makes a move
	time.Sleep(2 * time.Second)

	availableCols := board.GetAvailableColumns()
	if len(availableCols) == 0 {
		return -1
	}

	opponent := 1
	if b.playerNum == 1 {
		opponent = 2
	}

	// 90% chance to play optimally (makes bot beatable)
	playOptimally := rand.Intn(100) < 90

	if playOptimally {
		// Priority 1: Win if possible
		for _, col := range availableCols {
			if b.isWinningMove(board, col, b.playerNum) {
				return col
			}
		}

		// Priority 2: Block opponent from winning
		for _, col := range availableCols {
			if b.isWinningMove(board, col, opponent) {
				return col
			}
		}

		// Priority 3: Look for potential winning setups (two in a row)
		bestCol := b.findBestSetup(board, availableCols)
		if bestCol != -1 {
			return bestCol
		}

		// Priority 4: Strategic move - prefer center columns
		centerCols := []int{3, 2, 4, 1, 5, 0, 6}
		for _, col := range centerCols {
			for _, availCol := range availableCols {
				if col == availCol {
					return col
				}
			}
		}
	}

	// 10% chance: Make a random move (intentional mistake)
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

// findBestSetup looks for moves that create two-in-a-row opportunities
func (b *Bot) findBestSetup(board *game.Board, availableCols []int) int {
	bestScore := -1
	bestCol := -1

	for _, col := range availableCols {
		clonedBoard := board.Clone()
		row, ok := clonedBoard.MakeMove(col, b.playerNum)
		if !ok {
			continue
		}

		score := b.evaluatePosition(clonedBoard, row, col, b.playerNum)
		if score > bestScore {
			bestScore = score
			bestCol = col
		}
	}

	return bestCol
}

// evaluatePosition gives a score to a position based on potential winning lines
func (b *Bot) evaluatePosition(board *game.Board, row, col, player int) int {
	score := 0
	grid := board.GetGrid()

	// Check horizontal
	count := 1
	// Left
	for c := col - 1; c >= 0 && grid[row][c] == player; c-- {
		count++
	}
	// Right
	for c := col + 1; c < 7 && grid[row][c] == player; c++ {
		count++
	}
	if count >= 2 {
		score += count * 10
	}

	// Check vertical
	count = 1
	// Down
	for r := row + 1; r < 6 && grid[r][col] == player; r++ {
		count++
	}
	if count >= 2 {
		score += count * 10
	}

	// Check diagonal (top-left to bottom-right)
	count = 1
	// Up-left
	for r, c := row-1, col-1; r >= 0 && c >= 0 && grid[r][c] == player; r, c = r-1, c-1 {
		count++
	}
	// Down-right
	for r, c := row+1, col+1; r < 6 && c < 7 && grid[r][c] == player; r, c = r+1, c+1 {
		count++
	}
	if count >= 2 {
		score += count * 10
	}

	// Check diagonal (top-right to bottom-left)
	count = 1
	// Up-right
	for r, c := row-1, col+1; r >= 0 && c < 7 && grid[r][c] == player; r, c = r-1, c+1 {
		count++
	}
	// Down-left
	for r, c := row+1, col-1; r < 6 && c >= 0 && grid[r][c] == player; r, c = r+1, c-1 {
		count++
	}
	if count >= 2 {
		score += count * 10
	}

	// Prefer center column
	if col == 3 {
		score += 5
	} else if col == 2 || col == 4 {
		score += 3
	}

	return score
}