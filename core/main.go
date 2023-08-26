package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/staugaard/app-os/core/internal/auth"
	"github.com/staugaard/app-os/core/internal/config"
	"github.com/staugaard/app-os/core/internal/graph"
	"github.com/staugaard/app-os/core/internal/pb"
	"github.com/staugaard/app-os/core/internal/users"
	"github.com/staugaard/app-os/core/middleware"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
	authServer := auth.NewServer(usersServer)

	mux := http.NewServeMux()

	resolver := graph.NewResolver(usersServer)
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})
	srv := handler.NewDefaultServer(schema)
	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		if errors.Is(e, graph.ErrAccessDenied) {
			err.Extensions = map[string]interface{}{
				"type": "ACCESS_DENIED",
			}
		}

		return err
	})

	mux.HandleFunc("/auth/login", authServer.Login)
	mux.HandleFunc("/auth/logout", authServer.Logout)
	mux.Handle("/api/core/graphiql", playground.Handler("GraphQL playground", "/api/core/graph"))
	mux.Handle("/api/core/graph", middleware.ToContext(srv.ServeHTTP))
	mux.HandleFunc("/api/core/graph/schema.graphql", func(w http.ResponseWriter, r *http.Request) {
		formatter.NewFormatter(w).FormatSchema(schema.Schema())
	})

	log.Printf("Listening HTTP at %v", 8081)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
