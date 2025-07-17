package pubsub

import (
	"app/internal/models"
	"context"
	"encoding/json"
	"log"

	"app/internal/database"
)

// publish sends a message into the Pub/Sub topic.
func Publish(msg models.Message) {
	data, _ := json.Marshal(msg)
	count, err := database.RDB.Publish(context.Background(), "topic", data).Result()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		log.Println("warning: no subscribers for topic")
	}
}
