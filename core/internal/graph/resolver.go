package graph

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/staugaard/app-os/core/internal/pb"
	"github.com/staugaard/app-os/core/middleware"
	"github.com/staugaard/app-os/core/types"
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
