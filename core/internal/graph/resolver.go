package graph

import (
	"context"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
	"github.com/staugaard/app-os/core/internal/pb"
	"github.com/staugaard/app-os/core/middleware"
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

func (r *Resolver) CurrentUserID(ctx context.Context) uint32 {
	return middleware.UserIDFromContext(ctx)
}

func (r *Resolver) CurrentUser(ctx context.Context) *pb.User {
	r.currentUserOnce.Do(func() {
		userID := r.CurrentUserID(ctx)
		if userID == 0 {
			return
		}

		response, err := r.UsersService.GetById(ctx, &pb.GetUserByIdRequest{
			Id: userID,
		})
		if err != nil {
			graphql.AddError(ctx, err)
			return
		}
		r.currentUser = response.User
	})
	return r.currentUser
}
