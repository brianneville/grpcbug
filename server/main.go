package main

import (
	"context"
	_ "embed"
	"log"
	"net"

	"github.com/brianneville/grpcbug/defaults"
	"github.com/brianneville/grpcbug/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Service struct {
	proto.UnsafeMockRPCServer
}

func (s Service) Get(ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	return &proto.GetResponse{
		B: []byte(defaults.BigResponse),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+defaults.DefaultPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	proto.RegisterMockRPCServer(grpcServer, &Service{})
	if err := grpcServer.Serve(lis); err != nil {
		return
	}
}
