package graph

import (
	"context"
	"sync"

	"github.com/Season-04/app-os/core/internal/pb"
	"github.com/Season-04/app-os/core/middleware"
	"github.com/Season-04/app-os/core/types"
	"github.com/pkg/errors"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UsersService    pb.UsersServiceServer
	currentUser     *pb.User
	currentUserOnce sync.Once
}

var ErrAccessDenied = errors.New("Access Denied")

func NewResolver(UsersService pb.UsersServiceServer) *Resolver {
	return &Resolver{
		UsersService:    UsersService,
		currentUserOnce: sync.Once{},
	}
}

func (r *Resolver) CurrentUser(ctx context.Context) *types.User {
	return middleware.UserFromContext(ctx)
}
