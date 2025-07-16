package pubsub2

import (
	"app/internal/models"
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

var (
	pubsubCli *pubsub.Client
	chatTopic *pubsub.Topic
	sub       *pubsub.Subscription
	projectID = "artful-logic-152702"
	topicID   = "message-pub"
)

// publish sends a message into the Pub/Sub topic.
func Publish(msg models.Message) {
	data, _ := json.Marshal(msg)
	chatTopic.Publish(context.Background(), &pubsub.Message{Data: data})
}
