package rpc

import (
	pb "ghasedak-pubsub/api/proto/src"
	"ghasedak-pubsub/api/rpc/handler"
	"ghasedak-pubsub/pkg"
	"google.golang.org/grpc"
	"net"
)

type Grpc struct {
}

func NewGrpc() *Grpc {
	return &Grpc{}
}

func (g *Grpc) Initialize(address string) {
	pkg.GetLogger().Info("Listening grpc on", address)
	go g.startGrpcServer(address)
}

func (*Grpc) startGrpcServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		pkg.GetLogger().Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPublisherServer(s, &handler.PublisherServer{})
	pb.RegisterSubscriberServer(s, &handler.SubscriberServer{})
	if err := s.Serve(lis); err != nil {
		pkg.GetLogger().Fatalf("failed to serve: %v", err)
	}
}
