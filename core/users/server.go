package users

import (
	"context"

	"github.com/staugaard/app-os/core/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUsersServiceServer
}

func (s *Server) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}

func (s *Server) GetById(context.Context, *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetById not implemented")
}

var _ pb.UsersServiceServer = &Server{}
