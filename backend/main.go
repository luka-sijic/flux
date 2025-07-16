package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string            `json:"type"`
	Username string            `json:"user,omitempty"`
	Content  string            `json:"content,omitempty"`
	Users    map[string]string `json:"users,omitempty"`
	Log      []string          `json:"log,omitempty"`
}

type Client struct {
	conn     *websocket.Conn
	username string
	lastPong time.Time
}

var (
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	mu        sync.Mutex
	clients   = make(map[*Client]struct{})
	messages  = make([]string, 0, 100)
	pubsubCli *pubsub.Client
	chatTopic *pubsub.Topic
	sub       *pubsub.Subscription
	projectID = "artful-logic-152702"
	topicID   = "message-pub"
)

func main() {
	// 1) Root context that cancels on SIGINT/SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2) Pub/Sub setup
	var err error
	pubsubCli, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	chatTopic = pubsubCli.Topic(topicID)

	// 3) Ensure a subscription exists for *this* instance
	hostname, _ := os.Hostname()
	subID := fmt.Sprintf("chat-sub-%s", hostname)
	sub = ensureSubscription(ctx, subID, chatTopic)

	// 4) Start the background loops
	go subscribeLoop(ctx)
	go pingLoop(ctx)

	// 5) HTTP server
	http.HandleFunc("/ws", wsHandler)
	srv := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("WebSocket server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// 6) Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received, cleaning up...")

	// 7) Gracefully stop HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

// ensureSubscription creates the subscription if it doesn’t exist.
func ensureSubscription(ctx context.Context, subID string, topic *pubsub.Topic) *pubsub.Subscription {
	sub := pubsubCli.Subscription(subID)
	exists, err := sub.Exists(ctx)
	if err != nil {
		log.Fatalf("Subscription.Exists: %v", err)
	}
	if !exists {
		sub, err = pubsubCli.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
			Topic:             topic,
			AckDeadline:       20 * time.Second,
			RetentionDuration: 24 * time.Hour,
		})
		if err != nil {
			log.Fatalf("CreateSubscription: %v", err)
		}
		log.Printf("Created subscription %q", subID)
	} else {
		log.Printf("Reusing subscription %q", subID)
	}
	return sub
}

// subscribeLoop pulls messages from Pub/Sub and broadcasts them.
func subscribeLoop(ctx context.Context) {
	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		var msg Message
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
			c.conn.WriteJSON(msg)
		}
		mu.Unlock()
		m.Ack()
	})
	if err != nil && ctx.Err() == nil {
		log.Fatalf("sub.Receive: %v", err)
	}
}

// pingLoop sends pings and detects timed‑out clients.
func pingLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			var toRemove []*Client

			mu.Lock()
			// 1) Send ping
			for c := range clients {
				c.conn.WriteJSON(Message{Type: "ping"})
			}
			// 2) Collect timed‑out
			for c := range clients {
				if now.Sub(c.lastPong) > 60*time.Second {
					toRemove = append(toRemove, c)
				}
			}
			// 3) Remove them
			for _, c := range toRemove {
				delete(clients, c)
				c.conn.Close()
			}
			// 4) Broadcast updated user list
			broadcastUsers()
			mu.Unlock()
		}
	}
}

// wsHandler upgrades to WS, replays log, registers client, then handles messages.
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{conn: conn, lastPong: time.Now()}

	// 1) Register client and replay chat history
	mu.Lock()
	clients[client] = struct{}{}
	backlog := make([]string, len(messages))
	copy(backlog, messages)
	mu.Unlock()

	conn.WriteJSON(Message{Type: "log", Log: backlog})
	broadcastUsers()

	// 2) Read loop
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			break
		}
		switch msg.Type {
		case "join":
			client.username = msg.Username
			publish(msg)

		case "chat":
			publish(msg)

		case "pong":
			client.lastPong = time.Now()
		}
	}

	// 3) Cleanup on disconnect
	mu.Lock()
	delete(clients, client)
	mu.Unlock()
	broadcastUsers()
	conn.Close()
}

// broadcastUsers compiles the active user list and broadcasts it.
func broadcastUsers() {
	users := make(map[string]string)
	for c := range clients {
		users[c.username] = "active"
	}
	for c := range clients {
		c.conn.WriteJSON(Message{Type: "users", Users: users})
	}
}

// publish sends a message into the Pub/Sub topic.
func publish(msg Message) {
	data, _ := json.Marshal(msg)
	chatTopic.Publish(context.Background(), &pubsub.Message{Data: data})
}
