package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/docker/docker/client"
	"github.com/staugaard/app-os/core/internal/auth"
	"github.com/staugaard/app-os/core/internal/config"
	"github.com/staugaard/app-os/core/internal/graph"
	"github.com/staugaard/app-os/core/internal/pb"
	"github.com/staugaard/app-os/core/internal/users"
	"google.golang.org/grpc"
)

func main() {
	mode := "server"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	if mode == "server" {
		usersServer := users.NewServer()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			runGRPC(usersServer)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			auth.RunHTTPServer(usersServer)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			runHTTPServer(usersServer)
		}()

		wg.Wait()
	} else if mode == "run" {
		runCFG()
	}
}

func runCFG() {
	configDirectory := os.Getenv("APP_OS_CONFIG_DIRECTORY")
	if configDirectory == "" {
		configDirectory = "/config"
	}

	cfg := config.NewConfig(configDirectory)
	err := cfg.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatalf("failed to connect to docker: %v", err)
	}

	err = cfg.Run(context.Background(), dockerClient)
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}

func runGRPC(usersServer pb.UsersServiceServer) {
	s := grpc.NewServer()

	pb.RegisterUsersServiceServer(s, usersServer)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Listening gRPC at %s", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func runHTTPServer(usersServer pb.UsersServiceServer) {
	mux := http.NewServeMux()

	resolver := &graph.Resolver{
		UsersService: usersServer,
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	mux.Handle("/api/core/graphiql", playground.Handler("GraphQL playground", "/api/core/graph"))
	mux.Handle("/api/core/graph", srv)

	log.Printf("Listening HTTP at %v", 8081)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
