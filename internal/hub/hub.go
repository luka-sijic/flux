package hub

import (
	"sync"

	"github.com/luka-sijic/flux/internal/models"
)

type Hub struct {
	clients    map[string]map[*models.Client]struct{}
	broadcast  chan deliverReq
	register   chan *models.Client
	unregister chan *models.Client
	mu         sync.RWMutex
}

type deliverReq struct {
	user string
	msg  []byte
}

func New() *Hub {
	return &Hub{
		clients:    make(map[string]map[*models.Client]struct{}),
		broadcast:  make(chan deliverReq, 1024),
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
func (h *Hub) Broadcast(user string, msg []byte) {
	if user == "" || len(msg) == 0 {
		return
	}
	h.broadcast <- deliverReq{user: user, msg: msg}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if h.clients[c.Username] == nil {
				h.clients[c.Username] = make(map[*models.Client]struct{})
			}
			h.clients[c.Username][c] = struct{}{}
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if set, ok := h.clients[c.Username]; ok {
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
