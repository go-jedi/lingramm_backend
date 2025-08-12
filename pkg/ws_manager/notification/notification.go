package notification

import (
	"context"
	"sync"

	"github.com/fasthttp/websocket"
)

type ConnectionEntry struct {
	Connection *websocket.Conn
	Mu         sync.Mutex
	Cancel     context.CancelFunc
}

type Hub struct {
	mu          sync.RWMutex
	connections map[string]*ConnectionEntry
}

func New() *Hub {
	return &Hub{
		connections: make(map[string]*ConnectionEntry),
	}
}

// Set connection by telegram id.
func (h *Hub) Set(telegramID string, ce *ConnectionEntry) {
	h.mu.Lock()
	h.connections[telegramID] = ce
	h.mu.Unlock()
}

// Get connection by telegram id.
func (h *Hub) Get(telegramID string) (*ConnectionEntry, bool) {
	h.mu.RLock()
	ce, ok := h.connections[telegramID]
	h.mu.RUnlock()

	return ce, ok
}

// Delete connection by telegram id.
func (h *Hub) Delete(telegramID string) {
	h.mu.Lock()
	delete(h.connections, telegramID)
	h.mu.Unlock()
}
