package graph

import "github.com/staugaard/app-os/core/internal/pb"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UsersService pb.UsersServiceServer
}
