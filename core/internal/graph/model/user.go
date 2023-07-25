package model

import (
	"strconv"

	"github.com/staugaard/app-os/core/internal/pb"
)

type User struct {
	*pb.User
}

func (u *User) ID() string {
	return strconv.FormatUint(uint64(u.Id), 10)
}
