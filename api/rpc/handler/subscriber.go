package handler

import (
	"context"
	pb "ghasedak-pubsub/api/proto/src"
	"ghasedak-pubsub/pkg/pubsub"
	google_protobuf2 "github.com/golang/protobuf/ptypes/empty"
)

type SubscriberServer struct{}

func (s *SubscriberServer) CreateSubscription(ctx context.Context, req *pb.Subscription) (*pb.Subscription, error) {
	pulsarClient := pubsub.GetPulsar()
	err := pulsarClient.Subscribe(req.Name, req.Topic)
	if err != nil {
		return nil, err
	}
	return &pb.Subscription{Name: req.Name, Topic: req.Topic}, nil
}

func (s *SubscriberServer) Acknowledge(ctx context.Context, req *pb.AcknowledgeRequest) (*google_protobuf2.Empty, error) {
	panic("implement me")
}

func (s *SubscriberServer) Pull(ctx context.Context, req *pb.PullRequest) (*pb.PullResponse, error) {
	pulsarClient := pubsub.GetPulsar()
	r, err := pulsarClient.Receive(req.Subscription)
	if err != nil {
		return nil, err
	}
	msg := pb.ReceivedMessage{Message: &pb.PubSubMessage{Data: r.Value}}
	messages := []*pb.ReceivedMessage{&msg}
	return &pb.PullResponse{ReceivedMessages: messages}, nil
}

func (s *SubscriberServer) StreamingPull(stream pb.Subscriber_StreamingPullServer) error {
	panic("implement me")
}

func (s *SubscriberServer) Seek(ctx context.Context, req *pb.SeekRequest) (*pb.SeekResponse, error) {
	panic("implement me")
}
