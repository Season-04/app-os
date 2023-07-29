package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.35

import (
	"context"

	"github.com/staugaard/app-os/core/internal/graph/model"
	"github.com/staugaard/app-os/core/internal/pb"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input *model.CreateUserInput) (*model.User, error) {
	request := &pb.CreateUserRequest{
		Name:         input.Name,
		EmailAddress: input.EmailAddress,
		Password:     input.Password,
	}

	response, err := r.UsersService.CreateUser(ctx, request)
	if err != nil {
		return nil, err
	}

	return &model.User{User: response.User}, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	response, err := r.UsersService.List(ctx, &pb.ListRequest{})
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, len(response.Users))
	for i, user := range response.Users {
		users[i] = &model.User{User: user}
	}
	return users, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }