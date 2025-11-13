package models

import "time"

type Player struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Piece    int    `json:"piece" bson:"piece"` // 1 or 2
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
	ID          string      `json:"id" bson:"_id"`
	Player1     *Player     `json:"player1" bson:"player1"`
	Player2     *Player     `json:"player2" bson:"player2"`
	Board       [][]int     `json:"board" bson:"board"`
	CurrentTurn int         `json:"current_turn" bson:"current_turn"`
	Status      GameStatus  `json:"status" bson:"status"`
	Winner      *Player     `json:"winner,omitempty" bson:"winner,omitempty"`
	IsBot       bool        `json:"is_bot" bson:"is_bot"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updated_at"`
	FinishedAt  *time.Time  `json:"finished_at,omitempty" bson:"finished_at,omitempty"`
}

type Move struct {
	GameID string `json:"game_id" bson:"game_id"`
	Column int    `json:"column" bson:"column"`
	Player int    `json:"player" bson:"player"`
}

type GameEvent struct {
	Type      string      `json:"type" bson:"type"`
	GameID    string      `json:"game_id" bson:"game_id"`
	Data      interface{} `json:"data" bson:"data"`
	Timestamp time.Time   `json:"timestamp" bson:"timestamp"`
}

type LeaderboardEntry struct {
	Username   string `json:"username" bson:"username"`
	Wins       int    `json:"wins" bson:"wins"`
	Losses     int    `json:"losses" bson:"losses"`
	Draws      int    `json:"draws" bson:"draws"`
	TotalGames int    `json:"total_games" bson:"total_games"`
}