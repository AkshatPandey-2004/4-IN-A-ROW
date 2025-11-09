package websocket

import (
    "sync"
)

type Hub struct {
    clients    map[string]*Client
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]*Client),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.id] = client
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client.id]; ok {
                delete(h.clients, client.id)
                close(client.send)
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            h.mu.RLock()
            for _, client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client.id)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) GetClient(id string) *Client {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return h.clients[id]
}

func (h *Hub) SendToClient(clientID string, message []byte) {
    h.mu.RLock()
    client, ok := h.clients[clientID]
    h.mu.RUnlock()
    
    if ok {
        select {
        case client.send <- message:
        default:
        }
    }
}