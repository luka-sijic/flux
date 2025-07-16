package pubsub2

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

// ensureSubscription creates the subscription if it doesn’t exist.
func EnsureSubscription(ctx context.Context, pubsubCli *pubsub.Client, subID string, topic *pubsub.Topic) *pubsub.Subscription {
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
