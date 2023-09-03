package graph

import (
	"context"

	"github.com/Season-04/appos/core/internal/pb"
	"github.com/Season-04/appos/core/middleware"
	"github.com/Season-04/appos/core/types"
	"github.com/pkg/errors"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UsersService        pb.UsersServiceServer
	ApplicationsService pb.ApplicationsServiceServer
}

var ErrAccessDenied = errors.New("Access Denied")

func NewResolver(
	UsersService pb.UsersServiceServer,
	ApplicationsService pb.ApplicationsServiceServer,
) *Resolver {
	return &Resolver{
		UsersService:        UsersService,
		ApplicationsService: ApplicationsService,
	}
}

func (r *Resolver) CurrentUser(ctx context.Context) *types.User {
	return middleware.UserFromContext(ctx)
}
