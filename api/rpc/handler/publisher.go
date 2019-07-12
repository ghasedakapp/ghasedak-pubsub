package handler

import (
	"context"
	pb "ghasedak-pubsub/api/proto/src"
)

type PublisherServer struct{}

func (p *PublisherServer) CreateTopic(ctx context.Context, req *pb.Topic) (*pb.Topic, error) {
	panic("implement me")
}

func (p *PublisherServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	panic("implement me")
}
