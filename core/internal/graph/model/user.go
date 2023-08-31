package model

import (
	"strconv"

	"github.com/Season-04/app-os/core/internal/pb"
)

var RoleToProtobuf = map[UserRole]pb.UserRole{
	UserRoleAdmin: pb.UserRole_USER_ROLE_ADMIN,
	UserRoleUser:  pb.UserRole_USER_ROLE_USER,
}

var RoleFromProtobuf = map[pb.UserRole]UserRole{
	pb.UserRole_USER_ROLE_ADMIN: UserRoleAdmin,
	pb.UserRole_USER_ROLE_USER:  UserRoleUser,
}

type User struct {
	*pb.User
}

func (u *User) ID() string {
	return strconv.FormatUint(uint64(u.Id), 10)
}
