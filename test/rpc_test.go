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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	r, err := PubClient.CreateTopic(ctx, &pb.Topic{Name: "topic"})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, r.Name, "topic")
}

func TestPublishMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	topic := fmt.Sprintf("topic-%d", rand.Int31())
	textMessage := "salam"
	publishTextMessage(t, ctx, textMessage, topic)
}

func TestGrpcSubscriber(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	topic := fmt.Sprintf("topic-%d", rand.Int31())
	subscriptionName := fmt.Sprintf("subscription-%d", rand.Int31())
	textMessage := "salam"

	_, err := SubClient.CreateSubscription(ctx, &pb.Subscription{Name: subscriptionName, Topic: topic})
	if err != nil {
		t.Fatal(err)
	}

	msg := publishTextMessage(t, ctx, textMessage, topic)

	time.Sleep(1 * time.Second)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)

	r, err := SubClient.Pull(ctx, &pb.PullRequest{Subscription: subscriptionName})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(r.ReceivedMessages[0].Message.Data), textMessage)
	assert.Equal(t, *r.ReceivedMessages[0].Message, *msg)

}

func TestGrpcGetLastMessage(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	topic := fmt.Sprintf("topic-%d", rand.Int31())
	subscriptionName := fmt.Sprintf("subscription-%d", rand.Int31())
	textMessage1 := "salam1"
	textMessage2 := "salam2"

	publishTextMessage(t, ctx, textMessage1, topic)

	_, err := SubClient.CreateSubscription(ctx, &pb.Subscription{Name: subscriptionName, Topic: topic})
	if err != nil {
		t.Fatal(err)
	}

	publishTextMessage(t, ctx, textMessage2, topic)
	r, err := SubClient.Pull(ctx, &pb.PullRequest{Subscription: subscriptionName})
	assert.Equal(t, string(r.ReceivedMessages[0].Message.Data), textMessage2)

	_, err = SubClient.Pull(ctx, &pb.PullRequest{Subscription: subscriptionName})
	assert.Equal(t, err.Error(), "rpc error: code = DeadlineExceeded desc = context deadline exceeded")
}

func publishTextMessage(t *testing.T, ctx context.Context, body string, topic string) *pb.PubSubMessage {
	_, err := PubClient.CreateTopic(ctx, &pb.Topic{Name: topic})
	if err != nil {
		t.Fatal(err)
	}
	msg := &pb.PubSubMessage{Data: []byte(body)}
	_, err = PubClient.Publish(ctx, &pb.PublishRequest{Topic: topic, Messages: []*pb.PubSubMessage{msg}})
	if err != nil {
		t.Fatal(err)
	}
	return msg

}
