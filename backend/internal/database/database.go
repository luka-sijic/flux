package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file", err)
	}
}

func Connect() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASS")

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // no password if empty
		DB:       0,             // default DB
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v\n", err)
	}
	log.Println("Redis connected")
}

func Close() {
	if RDB != nil {
		if err := RDB.Close(); err != nil {
			log.Printf("Error closing redis connection: %v\n", err)
		}
	}
}
