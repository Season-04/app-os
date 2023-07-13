package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/staugaard/app-os/core/internal/config"
	"github.com/staugaard/app-os/core/internal/pb"
	"github.com/staugaard/app-os/core/internal/runner"
	"github.com/staugaard/app-os/core/internal/users"
	"google.golang.org/grpc"
)

func main() {
	configDirectory := os.Getenv("APP_OS_CONFIG_DIRECTORY")
	if configDirectory == "" {
		configDirectory = "/config"
	}

	cfg, err := config.Load(configDirectory)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	err = runner.Run(context.Background(), *cfg)
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}

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
