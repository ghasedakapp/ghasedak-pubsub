package rpc

import (
	pb "ghasedak-pubsub/api/proto/src"
	"ghasedak-pubsub/api/rpc/handler"
	"ghasedak-pubsub/internal/pubsub"
	"ghasedak-pubsub/pkg"
	"google.golang.org/grpc"
	"net"
)

type Grpc struct {
	log *pkg.Logger
}

func NewGrpc(log *pkg.Logger, adapter *pubsub.Adapter, address string) *Grpc {
	g := &Grpc{}
	log.Info("Listening grpc on", address)
	go g.startGrpcServer(address, adapter)
	return g
}

func (g *Grpc) startGrpcServer(address string, adapter *pubsub.Adapter) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		g.log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPublisherServer(s, handler.NewPublisherServer(adapter))
	pb.RegisterSubscriberServer(s, handler.NewSubscriberServer(adapter))
	if err := s.Serve(lis); err != nil {
		g.log.Fatalf("failed to serve: %v", err)
	}
}
