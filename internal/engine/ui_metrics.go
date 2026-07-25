package engine

import (
	"github.com/gorilla/websocket"
	"sync"
)

type UIMetricsHub struct {
	mu      sync.RWMutex
	clients map[string]map[*websocket.Conn]bool
}

var GlobalUIMetricsHub = &UIMetricsHub{
	clients: make(map[string]map[*websocket.Conn]bool),
}

func (h *UIMetricsHub) AddClient(serverID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[serverID] == nil {
		h.clients[serverID] = make(map[*websocket.Conn]bool)
	}
	h.clients[serverID][conn] = true
}

func (h *UIMetricsHub) RemoveClient(serverID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[serverID] != nil {
		delete(h.clients[serverID], conn)
	}
}

func (h *UIMetricsHub) Broadcast(serverID string, metrics []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for conn := range h.clients[serverID] {
		// Non-blocking write or fire-and-forget
		go conn.WriteMessage(websocket.TextMessage, metrics)
	}
}
