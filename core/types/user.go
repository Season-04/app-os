package types

import "time"

type UserRole string

var (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type User struct {
	ID           uint32     `json:"id"`
	Name         string     `json:"name"`
	EmailAddress string     `json:"email_address"`
	Role         UserRole   `json:"role"`
	LastSeenAt   *time.Time `json:"last_seen_at"`
}
