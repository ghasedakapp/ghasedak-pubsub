package test

import (
	"context"
	pb "ghasedak-pubsub/api/proto/src"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGrpcSubscriber(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := PubClient.CreateTopic(ctx, &pb.Topic{Name: "topic"})
	if err != nil {
		t.Fatal(err)
	}

	msg := &pb.PubSubMessage{Data: []byte("salam")}
	_, err = PubClient.Publish(ctx, &pb.PublishRequest{Topic: "topic", Messages: []*pb.PubSubMessage{msg}})
	if err != nil {
		t.Fatal(err)
	}

	_, err = SubClient.CreateSubscription(ctx, &pb.Subscription{Name: "test", Topic: "topic"})
	if err != nil {
		t.Fatal(err)
	}

	r, err := SubClient.Pull(ctx, &pb.PullRequest{Subscription: "test"})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(r.ReceivedMessages[0].Message.Data), "salam")
	assert.Equal(t, *r.ReceivedMessages[0].Message, *msg)

}
