CREATE TABLE IF NOT EXISTS players (
    id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS games (
    id VARCHAR(255) PRIMARY KEY,
    player1_id VARCHAR(255) REFERENCES players(id),
    player2_id VARCHAR(255) REFERENCES players(id),
    winner_id VARCHAR(255) REFERENCES players(id),
    is_bot BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS game_stats (
    username VARCHAR(255) PRIMARY KEY,
    wins INT DEFAULT 0,
    losses INT DEFAULT 0,
    draws INT DEFAULT 0,
    total_games INT DEFAULT 0
);

CREATE INDEX idx_games_player1 ON games(player1_id);
CREATE INDEX idx_games_player2 ON games(player2_id);
CREATE INDEX idx_games_status ON games(status);