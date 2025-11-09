package game

const (
    Rows = 6
    Cols = 7
)

type Board struct {
    grid [][]int
}

func NewBoard() *Board {
    grid := make([][]int, Rows)
    for i := range grid {
        grid[i] = make([]int, Cols)
    }
    return &Board{grid: grid}
}

func (b *Board) GetGrid() [][]int {
    return b.grid
}

func (b *Board) IsValidMove(col int) bool {
    if col < 0 || col >= Cols {
        return false
    }
    return b.grid[0][col] == 0
}

func (b *Board) MakeMove(col int, player int) (int, bool) {
    if !b.IsValidMove(col) {
        return -1, false
    }

    for row := Rows - 1; row >= 0; row-- {
        if b.grid[row][col] == 0 {
            b.grid[row][col] = player
            return row, true
        }
    }
    return -1, false
}

func (b *Board) CheckWin(row, col, player int) bool {
    // Check horizontal
    if b.checkDirection(row, col, 0, 1, player) {
        return true
    }
    // Check vertical
    if b.checkDirection(row, col, 1, 0, player) {
        return true
    }
    // Check diagonal (top-left to bottom-right)
    if b.checkDirection(row, col, 1, 1, player) {
        return true
    }
    // Check diagonal (bottom-left to top-right)
    if b.checkDirection(row, col, -1, 1, player) {
        return true
    }
    return false
}

func (b *Board) checkDirection(row, col, dRow, dCol, player int) bool {
    count := 1

    // Check positive direction
    r, c := row+dRow, col+dCol
    for r >= 0 && r < Rows && c >= 0 && c < Cols && b.grid[r][c] == player {
        count++
        r += dRow
        c += dCol
    }

    // Check negative direction
    r, c = row-dRow, col-dCol
    for r >= 0 && r < Rows && c >= 0 && c < Cols && b.grid[r][c] == player {
        count++
        r -= dRow
        c -= dCol
    }

    return count >= 4
}

func (b *Board) IsFull() bool {
    for col := 0; col < Cols; col++ {
        if b.grid[0][col] == 0 {
            return false
        }
    }
    return true
}

func (b *Board) GetAvailableColumns() []int {
    available := []int{}
    for col := 0; col < Cols; col++ {
        if b.IsValidMove(col) {
            available = append(available, col)
        }
    }
    return available
}

func (b *Board) Clone() *Board {
    newGrid := make([][]int, Rows)
    for i := range b.grid {
        newGrid[i] = make([]int, Cols)
        copy(newGrid[i], b.grid[i])
    }
    return &Board{grid: newGrid}
}