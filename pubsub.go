package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()

	projectID := "artful-logic-152702"
	topicID := "message-pub"
	subID := "message-sub"

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	topic := client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("Topic.Exists: %v", err)
	}
	if !exists {
		topic, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			log.Fatalf("CreateTopic: %v", err)
		}
		fmt.Println("Created topic:", topicID)
	}

	sub := client.Subscription(subID)
	exists, err = sub.Exists(ctx)
	if err != nil {
		log.Fatalf("Subscription.Exists: %v", err)
	}
	if !exists {
		sub, err = client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			log.Fatalf("CreateSubscription: %v", err)
		}
		fmt.Println("Created subscription:", subID)
	}

	// 4) Publish a message
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(`{"user":"alice","content":"Hello, world!"}`),
		Attributes: map[string]string{
			"origin": "go-sample",
		},
	})

	// Block until the publish is acknowledged
	id, err := result.Get(ctx)
	if err != nil {
		log.Fatalf("Publish.Get: %v", err)
	}
	fmt.Println("Published message ID:", id)

	// 5) Receive messages
	fmt.Println("Pulling messages (will exit after 5s)…")
	cctx, cancel := context.WithCancel(ctx)
	// Stop after 5s
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	err = sub.Receive(cctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Printf("Got message:\n  Data: %s\n  Attributes: %v\n", string(m.Data), m.Attributes)
		m.Ack() // acknowledge so it won’t be resent
	})
	if err != nil && err != context.Canceled {
		log.Fatalf("Receive: %v", err)
	}

	// 6) List all topics
	fmt.Println("\nTopics in project:")
	it := client.Topics(ctx)
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Topics.Next: %v", err)
		}
		fmt.Println(" -", t.ID())
	}

	// 7) List all subscriptions
	fmt.Println("\nSubscriptions in project:")
	sit := client.Subscriptions(ctx)
	for {
		s, err := sit.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Subscriptions.Next: %v", err)
		}
		fmt.Println(" -", s.ID())
	}

	os.Exit(0)
}
