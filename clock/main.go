package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Season-04/app-os/clock/pb"
	"github.com/Season-04/app-os/core/middleware"
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
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		runGRPC()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runHTTPServer()
	}()

	wg.Wait()
}

func runGRPC() {
	s := grpc.NewServer()
	pb.RegisterClockServiceServer(s, &ClockServer{})

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

func runHTTPServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/clock", middleware.ToContext(func(w http.ResponseWriter, r *http.Request) {
		user := middleware.UserFromContext(r.Context())
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(fmt.Sprintf("Hi %v, it is %s", user.Name, time.Now().String())))
	}))

	log.Printf("Listening HTTP at %v", 8081)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
