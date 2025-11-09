package game

import (
    "sync"
)

type Manager struct {
    games map[string]*GameInstance
    mu    sync.RWMutex
}

func NewManager() *Manager {
    return &Manager{
        games: make(map[string]*GameInstance),
    }
}

func (m *Manager) AddGame(game *GameInstance) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.games[game.ID] = game
}

func (m *Manager) GetGame(gameID string) (*GameInstance, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    game, ok := m.games[gameID]
    return game, ok
}

func (m *Manager) RemoveGame(gameID string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.games, gameID)
}

func (m *Manager) GetAllActiveGames() []*GameInstance {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    games := make([]*GameInstance, 0, len(m.games))
    for _, game := range m.games {
        games = append(games, game)
    }
    return games
}