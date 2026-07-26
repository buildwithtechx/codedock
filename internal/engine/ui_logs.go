package engine

import (
	"sync"

	"github.com/gorilla/websocket"
)

type UILogStreamHub struct {
	mu      sync.RWMutex
	clients map[string]map[*websocket.Conn]bool
}

var GlobalUILogStreamHub = &UILogStreamHub{
	clients: make(map[string]map[*websocket.Conn]bool),
}

func (h *UILogStreamHub) AddClient(key string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[key] == nil {
		h.clients[key] = make(map[*websocket.Conn]bool)
	}
	h.clients[key][conn] = true
}

func (h *UILogStreamHub) RemoveClient(key string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[key] != nil {
		delete(h.clients[key], conn)
	}
}

func (h *UILogStreamHub) Broadcast(key string, logLine []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for conn := range h.clients[key] {
		go conn.WriteMessage(websocket.TextMessage, logLine)
	}
}
