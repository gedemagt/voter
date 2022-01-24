package domain

import "github.com/google/uuid"

type Role string

const (
	SuperAdmin  Role = "SUPER_ADMIN"
	PollAdmin   Role = "POLL_ADMIN"
	RegularUser Role = "REGULAR_USER"
	Temporary   Role = "TEMPORARY"
)

type User struct {
	UUID  uuid.UUID
	Name  string
	Email string
	Role  Role
}

type UserRepository interface {
	GetUsers() ([]User, error)
	GetUserByUUID(uuid uuid.UUID) (User, error)
	GetUserByEmail(email string) (User, error)
	CreateUser(user User) error
	DeleteUser(user User) error
}
