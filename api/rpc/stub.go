package rpc

import (
	pb "ghasedak-pubsub/api/proto/src"
	"ghasedak-pubsub/api/rpc/handler"
	"ghasedak-pubsub/pkg"
	"google.golang.org/grpc"
	"net"
)

func InitGrpc(address string) {
	pkg.Logger.Info("Listening grpc on", address)
	go startGrpcServer(address)
}

func startGrpcServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		pkg.Logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPublisherServer(s, &handler.PublisherServer{})
	pb.RegisterSubscriberServer(s, &handler.SubscriberServer{})
	if err := s.Serve(lis); err != nil {
		pkg.Logger.Fatalf("failed to serve: %v", err)
	}
}
