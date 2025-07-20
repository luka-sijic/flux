package transport

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/luka-sijic/flux/internal/hub"
	"github.com/luka-sijic/flux/internal/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	Hub *hub.Hub
}

func New(h *hub.Hub) http.Handler { return &Handler{Hub: h} }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	client := &models.Client{
		Conn:     conn,
		Send:     make(chan []byte, 256),
		LastPong: time.Now(),
	}

	h.Hub.Register(client)

	go readPump(h.Hub, client)
	go writePump(client)
}

const (
	pingInterval = 30 * time.Second
	pongWait     = 60 * time.Second
)

func readPump(h *hub.Hub, c *models.Client) {
	defer func() {
		h.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.LastPong = time.Now()
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		h.Broadcast(msg)
	}
}

func writePump(c *models.Client) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = c.Conn.WriteMessage(websocket.TextMessage, msg)

		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
