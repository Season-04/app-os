package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

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

	mux.HandleFunc("/clock", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-AppOS-User-ID")
		w.Write([]byte(fmt.Sprintf("User %v, it is %s", userID, time.Now().String())))
	})

	log.Printf("Listening HTTP at %v", 8081)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
