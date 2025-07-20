package server

/*
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/luka-sijic/flux/internal/database"
	"github.com/luka-sijic/flux/internal/models"

	"github.com/gorilla/websocket"
)

type Conn struct {
	websocket.Upgrader
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	mu       sync.Mutex
	clients  = make(map[*models.Client]struct{})
)
}

func PingLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			var toRemove []*models.Client

			mu.Lock()
			// 1) Send ping
			for c := range clients {
				c.Conn.WriteJSON(models.Message{Type: "ping"})
			}
			// 2) Collect timed‑out
			for c := range clients {
				if now.Sub(c.LastPong) > 60*time.Second {
					toRemove = append(toRemove, c)
				}
			}
			// 3) Remove them
			for _, c := range toRemove {
				delete(clients, c)
				c.Conn.Close()
			}
			// 4) Broadcast updated user list
			broadcastUsers()
			mu.Unlock()
		}
	}
}

func SubscribeLoop(ctx context.Context) {
	sub := database.RDB.Subscribe(ctx, "topic")
	defer sub.Close()

	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case m, ok := <-ch:
			if !ok {
				// channel closed
				return
			}

			var msg models.Message
			if err := json.Unmarshal([]byte(m.Payload), &msg); err != nil {
				log.Println("pubsub unmarshal:", err)
				continue
			}

			// build history & broadcast
			/*if msg.Type == "chat" {
				line := fmt.Sprintf("%s: %s", msg.Username, msg.Content)
				mu.Lock()
				//messages = append(messages, line)
				mu.Unlock()
			}

			// snapshot clients under RLock
			mu.Lock()
			conns := make([]*models.Client, 0, len(clients))
			for c := range clients {
				conns = append(conns, c)
			}
			mu.Unlock()

			// broadcast without holding the lock
			for _, c := range conns {
				if err := c.Conn.WriteJSON(msg); err != nil {
					log.Printf("ws write error (%v): %v", c, err)
				}
			}
		}
	}
}


broadcastUsers compiles the active user list and broadcasts it.

	func broadcastUserJoined(username string) {
		for c := range clients {
			c.Conn.WriteJSON(models.Message{Type: "user_joined", User: username})
		}
	}

func broadcastUsers() {
	users := make(map[string]string)
	for c := range clients {
		if c.Username != "" {
			users[c.Username] = "active"
		}
	}
	for c := range clients {
		c.Conn.WriteJSON(models.Message{Type: "users", Users: users})
	}
}
*/
