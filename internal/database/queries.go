package database

import (
    "database/sql"
    "github.com/AkshatPandey-2004/4-in-a-row/pkg/models"
)

func (db *DB) CreateOrGetPlayer(username string, playerID string) error {
    query := `INSERT INTO players (id, username) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING`
    _, err := db.Exec(query, playerID, username)
    return err
}

func (db *DB) SaveGame(game *models.Game) error {
    query := `
        INSERT INTO games (id, player1_id, player2_id, winner_id, is_bot, status, created_at, finished_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
    
    var player1ID, player2ID, winnerID sql.NullString
    
    if game.Player1 != nil {
        player1ID.String = game.Player1.ID
        player1ID.Valid = true
    }
    
    if game.Player2 != nil {
        player2ID.String = game.Player2.ID
        player2ID.Valid = true
    }
    
    if game.Winner != nil {
        winnerID.String = game.Winner.ID
        winnerID.Valid = true
    }
    
    _, err := db.Exec(query, game.ID, player1ID, player2ID, winnerID, 
        game.IsBot, game.Status, game.CreatedAt, game.FinishedAt)
    
    return err
}

func (db *DB) UpdateGameStats(game *models.Game) error {
    if game.Winner != nil {
        // Update winner stats
        _, err := db.Exec(`
            INSERT INTO game_stats (username, wins, total_games) VALUES ($1, 1, 1)
            ON CONFLICT (username) DO UPDATE SET wins = game_stats.wins + 1, total_games = game_stats.total_games + 1
        `, game.Winner.Username)
        if err != nil {
            return err
        }
        
        // Update loser stats
        loser := game.Player1
        if game.Winner.ID == game.Player1.ID {
            loser = game.Player2
        }
        
        if loser != nil {
            _, err = db.Exec(`
                INSERT INTO game_stats (username, losses, total_games) VALUES ($1, 1, 1)
                ON CONFLICT (username) DO UPDATE SET losses = game_stats.losses + 1, total_games = game_stats.total_games + 1
            `, loser.Username)
            return err
        }
    } else {
        // Draw - update both players
        if game.Player1 != nil {
            db.Exec(`
                INSERT INTO game_stats (username, draws, total_games) VALUES ($1, 1, 1)
                ON CONFLICT (username) DO UPDATE SET draws = game_stats.draws + 1, total_games = game_stats.total_games + 1
            `, game.Player1.Username)
        }
        
        if game.Player2 != nil {
            db.Exec(`
                INSERT INTO game_stats (username, draws, total_games) VALUES ($1, 1, 1)
                ON CONFLICT (username) DO UPDATE SET draws = game_stats.draws + 1, total_games = game_stats.total_games + 1
            `, game.Player2.Username)
        }
    }
    
    return nil
}

func (db *DB) GetLeaderboard(limit int) ([]models.LeaderboardEntry, error) {
    query := `
        SELECT username, wins, losses, draws, total_games 
        FROM game_stats 
        ORDER BY wins DESC, total_games DESC 
        LIMIT $1
    `
    
    var leaderboard []models.LeaderboardEntry
    err := db.Select(&leaderboard, query, limit)
    return leaderboard, err
}