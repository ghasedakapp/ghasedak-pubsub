package pubsub

import (
	"errors"
	"time"
)

var (
	SubscriptionNotFound = errors.New("subscription not found")
	TopicNotFound        = errors.New("topic not found")
)

type MessageId struct {
	id []byte
}

type Message struct {
	Value     []byte
	Timestamp time.Time
}

type PubSub interface {
	Publish(topic string, body []byte) (*MessageId, error)
	CreateProducer(topic string) error
	Subscribe(subscriptionName string, topic string) error
	Receive(subscriptionName string) (*Message, error)
	Ack(subscriptionName string, mid MessageId) error
	Close() error
}
