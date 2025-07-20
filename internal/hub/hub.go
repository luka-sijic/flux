package hub

import (
	"sync"

	"github.com/luka-sijic/flux/internal/models"
)

type Hub struct {
	clients    map[*models.Client]struct{}
	broadcast  chan []byte
	register   chan *models.Client
	unregister chan *models.Client
	mu         sync.RWMutex
}

func New() *Hub {
	return &Hub{
		clients:    make(map[*models.Client]struct{}),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *models.Client),
		unregister: make(chan *models.Client),
	}
}

func (h *Hub) Register(c *models.Client) {
	if c == nil {
		return
	}
	h.register <- c
}
func (h *Hub) Unregister(c *models.Client) {
	if c == nil {
		return
	}
	h.unregister <- c
}
func (h *Hub) Broadcast(msg []byte) {
	if len(msg) == 0 {
		return
	}
	h.broadcast <- msg
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.Send <- msg:
				default:
					close(c.Send)
					delete(h.clients, c)
				}
			}
			h.mu.RUnlock()
		}
	}
}
