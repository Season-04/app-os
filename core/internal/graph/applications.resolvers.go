package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"

	"github.com/Season-04/appos/core/internal/graph/model"
	"github.com/Season-04/appos/core/internal/pb"
)

// InstalledApplications is the resolver for the installedApplications field.
func (r *queryResolver) InstalledApplications(ctx context.Context) ([]*model.Application, error) {
	response, err := r.ApplicationsService.ListInstalledApplications(ctx, &pb.ListInstalledApplicationsRequest{})
	if err != nil {
		return nil, err
	}

	apps := make([]*model.Application, len(response.Applications))

	for i, a := range response.Applications {
		apps[i] = &model.Application{
			ID:   a.Id,
			Name: a.Name,
		}
	}

	return apps, nil
}