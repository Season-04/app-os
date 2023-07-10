package main

import (
	"context"
	"log"
	"net"

	"github.com/staugaard/app-os/clock/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ClockServer struct {
	pb.UnimplementedClockServiceServer
}

func (s *ClockServer) GetTime(ctx context.Context, r *pb.GetTimeRequest) (*pb.GetTimeResponse, error) {
	return &pb.GetTimeResponse{
		Time: timestamppb.Now(),
	}, nil
}

func main() {
	s := grpc.NewServer()
	pb.RegisterClockServiceServer(s, &ClockServer{})

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening at %s", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
