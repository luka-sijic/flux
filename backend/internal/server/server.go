package server

import (
	"app/internal/database"
	"app/internal/models"
	"app/internal/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	mu       sync.Mutex
	clients  = make(map[*models.Client]struct{})
	messages = make([]string, 0, 100)
)

// wsHandler upgrades to WS, replays log, registers client, then handles messages.
func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &models.Client{Conn: conn, LastPong: time.Now()}

	// 1) Register client and replay chat history
	mu.Lock()
	clients[client] = struct{}{}
	backlog := make([]string, len(messages))
	copy(backlog, messages)
	mu.Unlock()

	conn.WriteJSON(models.Message{Type: "log", Log: backlog})
	broadcastUsers()

	// 2) Read loop
	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			break
		}
		switch msg.Type {
		case "join":
			client.Username = msg.Username
			pubsub.Publish(msg)

		case "chat":
			pubsub.Publish(msg)

		case "pong":
			client.LastPong = time.Now()
		}
	}

	// 3) Cleanup on disconnect
	mu.Lock()
	delete(clients, client)
	mu.Unlock()
	broadcastUsers()
	conn.Close()
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
			if msg.Type == "chat" {
				line := fmt.Sprintf("%s: %s", msg.Username, msg.Content)
				mu.Lock()
				messages = append(messages, line)
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

/*
broadcastUsers compiles the active user list and broadcasts it.

	func broadcastUserJoined(username string) {
		for c := range clients {
			c.Conn.WriteJSON(models.Message{Type: "user_joined", User: username})
		}
	}
*/
func broadcastUsers() {
	users := make(map[string]string)
	fmt.Println(users)
	for c := range clients {
		if c.Username != "" {
			users[c.Username] = "active"
		}
	}
	for c := range clients {
		c.Conn.WriteJSON(models.Message{Type: "users", Users: users})
	}
}
