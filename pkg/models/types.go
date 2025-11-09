package models

import "time"

type Player struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Piece    int    `json:"piece"` // 1 or 2
}

type GameStatus string

const (
    StatusWaiting  GameStatus = "waiting"
    StatusPlaying  GameStatus = "playing"
    StatusFinished GameStatus = "finished"
)

type GameResult string

const (
    ResultWin  GameResult = "win"
    ResultLoss GameResult = "loss"
    ResultDraw GameResult = "draw"
)

type Game struct {
    ID          string      `json:"id"`
    Player1     *Player     `json:"player1"`
    Player2     *Player     `json:"player2"`
    Board       [][]int     `json:"board"`
    CurrentTurn int         `json:"current_turn"`
    Status      GameStatus  `json:"status"`
    Winner      *Player     `json:"winner,omitempty"`
    IsBot       bool        `json:"is_bot"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`
    FinishedAt  *time.Time  `json:"finished_at,omitempty"`
}

type Move struct {
    GameID string `json:"game_id"`
    Column int    `json:"column"`
    Player int    `json:"player"`
}

type GameEvent struct {
    Type      string      `json:"type"`
    GameID    string      `json:"game_id"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}

type LeaderboardEntry struct {
    Username  string `json:"username"`
    Wins      int    `json:"wins"`
    Losses    int    `json:"losses"`
    Draws     int    `json:"draws"`
    TotalGames int   `json:"total_games"`
}