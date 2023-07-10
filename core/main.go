package main

import (
	"log"
	"net"

	"github.com/staugaard/app-os/core/pb"
	"github.com/staugaard/app-os/core/users"
	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()
	pb.RegisterUsersServiceServer(s, &users.Server{})

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
