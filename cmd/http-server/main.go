package main

import (
	"log"

	"github.com/luka-sijic/flux/internal/server"
)

func main() {
	log.Println("users table is ready on all shards")
	server.Start()
}
