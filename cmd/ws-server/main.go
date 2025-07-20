package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luka-sijic/flux/internal/database"
	"github.com/luka-sijic/flux/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	database.Connect()
	defer database.Close()

	go server.SubscribeLoop(ctx)
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
