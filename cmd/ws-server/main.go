package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/luka-sijic/flux/internal/hub"
	"github.com/luka-sijic/flux/internal/transport"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	h := hub.New()
	go h.Run()

	//go server.SubscribeLoop(ctx)
	//go server.PingLoop(ctx)

	r := mux.NewRouter()
	r.Handle("/ws", transport.New(h)).Methods("GET")

	srv := &http.Server{
		Addr:    ":8085",
		Handler: r,
	}

	go func() {
		log.Println("WebSocket server listening on :8085")
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
