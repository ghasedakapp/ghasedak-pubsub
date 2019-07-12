package test

import (
	"context"
	"fmt"
	pb "ghasedak-pubsub/api/proto/src"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestCreateTopic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := PubClient.CreateTopic(ctx, &pb.Topic{Name: "topic"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, r.Name, "topic")
}

func TestGrpcSubscriber(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	topic := fmt.Sprintf("topic-%d", rand.Int31())
	subscriptionName := fmt.Sprintf("subscription-%d", rand.Int31())
	textMessage := "salam"

	_, err := SubClient.CreateSubscription(ctx, &pb.Subscription{Name: subscriptionName, Topic: topic})
	if err != nil {
		t.Fatal(err)
	}

	_, err = PubClient.CreateTopic(ctx, &pb.Topic{Name: topic})
	if err != nil {
		t.Fatal(err)
	}

	msg := &pb.PubSubMessage{Data: []byte(textMessage)}
	_, err = PubClient.Publish(ctx, &pb.PublishRequest{Topic: topic, Messages: []*pb.PubSubMessage{msg}})
	if err != nil {
		t.Fatal(err)
	}

	r, err := SubClient.Pull(ctx, &pb.PullRequest{Subscription: subscriptionName})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(r.ReceivedMessages[0].Message.Data), textMessage)
	assert.Equal(t, *r.ReceivedMessages[0].Message, *msg)

}
