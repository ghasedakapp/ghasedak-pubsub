package handler

import (
	"context"
	pb "ghasedak-pubsub/api/proto/src"
	pubsub2 "ghasedak-pubsub/internal/pubsub"
)

type PublisherServer struct {
	pubsub *pubsub2.Adapter
}

func NewPublisherServer(pbAdapter *pubsub2.Adapter) *PublisherServer {
	return &PublisherServer{
		pubsub: pbAdapter,
	}
}

func (p *PublisherServer) CreateTopic(ctx context.Context, req *pb.Topic) (*pb.Topic, error) {
	err := p.pubsub.CreateProducer(req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.Topic{Name: req.Name}, nil
}

func (p *PublisherServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	var ids []*pb.MessageId
	for _, m := range req.Messages {
		_, err := p.pubsub.Publish(req.Topic, m.Data)
		if err != nil {
			return nil, err
		}

		ids = append(ids, &pb.MessageId{})
	}
	return &pb.PublishResponse{MessageIds: ids}, nil
}
