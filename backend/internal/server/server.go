package server

import (
	"app/internal/models"
	"app/internal/pubsub2"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
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
			pubsub2.Publish(msg)

		case "chat":
			pubsub2.Publish(msg)

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

// subscribeLoop pulls messages from Pub/Sub and broadcasts them.
func SubscribeLoop(ctx context.Context, sub *pubsub.Subscription) {
	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		var msg models.Message
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			log.Println("pubsub unmarshal:", err)
			m.Ack()
			return
		}

		mu.Lock()
		// If it's a chat message, append to history
		if msg.Type == "chat" {
			line := fmt.Sprintf("%s: %s", msg.Username, msg.Content)
			messages = append(messages, line)
		}
		// Broadcast to all WS clients
		for c := range clients {
			c.Conn.WriteJSON(msg)
		}
		mu.Unlock()
		m.Ack()
	})
	if err != nil && ctx.Err() == nil {
		log.Fatalf("sub.Receive: %v", err)
	}
}

// broadcastUsers compiles the active user list and broadcasts it.
func broadcastUsers() {
	users := make(map[string]string)
	for c := range clients {
		users[c.Username] = "active"
	}
	for c := range clients {
		c.Conn.WriteJSON(models.Message{Type: "users", Users: users})
	}
}
