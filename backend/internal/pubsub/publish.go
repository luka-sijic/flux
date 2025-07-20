package pubsub

import (
	"app/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"app/internal/database"

	"github.com/redis/go-redis/v9"
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
	var key string
	if msg.Username > msg.User2 {
		key = fmt.Sprintf("conversation:%s:%s", msg.Username, msg.User2)
	} else {
		key = fmt.Sprintf("conversation:%s:%s", msg.User2, msg.Username)
	}
	fmt.Println("Key: ", key)
	_, err = database.RDB.XAdd(context.Background(), &redis.XAddArgs{
		Stream: key,
		Values: map[string]interface{}{
			"username": msg.Username,
			"message":  msg.Content,
		},
	}).Result()
	if err != nil {
		log.Println(err)
	}
}
