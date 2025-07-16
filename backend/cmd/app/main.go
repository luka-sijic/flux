package main

import (
	"app/internal/pubsub2"
	"app/internal/server"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
)

var (
	pubsubCli *pubsub.Client
	chatTopic *pubsub.Topic
	sub       *pubsub.Subscription
	projectID = "artful-logic-152702"
	topicID   = "message-pub"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var err error
	pubsubCli, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	chatTopic = pubsubCli.Topic(topicID)

	// 3) Ensure a subscription exists for *this* instance
	hostname, _ := os.Hostname()
	subID := fmt.Sprintf("chat-sub-%s", hostname)
	sub = pubsub2.EnsureSubscription(ctx, pubsubCli, subID, chatTopic)

	go server.SubscribeLoop(ctx, sub)
	go server.PingLoop(ctx)

	http.HandleFunc("/ws", server.WSHandler)
	srv := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("WebSocket server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, cleaning up...")

	// 7) Gracefully stop HTTP server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
