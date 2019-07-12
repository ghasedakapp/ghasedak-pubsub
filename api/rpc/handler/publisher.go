package handler

import (
	"context"
	pb "ghasedak-pubsub/api/proto/src"
	"ghasedak-pubsub/pkg/pubsub"
)

type PublisherServer struct{}

func (p *PublisherServer) CreateTopic(ctx context.Context, req *pb.Topic) (*pb.Topic, error) {
	pulsarClient := pubsub.GetPulsar()
	err := pulsarClient.CreateProducer(req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.Topic{Name: req.Name}, nil
}

func (p *PublisherServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	pulsarClient := pubsub.GetPulsar()
	var ids []*pb.MessageId
	for _, m := range req.Messages {
		_, err := pulsarClient.Publish(req.Topic, m.Data)
		if err != nil {
			return nil, err
		}

		ids = append(ids, &pb.MessageId{})
	}
	return &pb.PublishResponse{MessageIds: ids}, nil
}
