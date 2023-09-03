package applications

import (
	"context"

	"github.com/Season-04/appos/core/internal/config"
	"github.com/Season-04/appos/core/internal/pb"
)

type Server struct {
	pb.UnimplementedApplicationsServiceServer
	cfg *config.Config
}

func (s *Server) ListInstalledApplications(
	ctx context.Context,
	r *pb.ListInstalledApplicationsRequest,
) (*pb.ListInstalledApplicationsResponse, error) {
	apps := make([]*pb.Application, len(s.cfg.Manifests()))

	for i, m := range s.cfg.Manifests() {
		apps[i] = &pb.Application{
			Id:   m.ID,
			Name: m.Name,
		}
	}

	return &pb.ListInstalledApplicationsResponse{
		Applications: apps,
	}, nil
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

var _ pb.ApplicationsServiceServer = &Server{}
