package identity

import (
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("forbidden: insufficient permissions")
)

type User struct {
	ID           string    `json:"id"`
	RoleID       string    `json:"role_id"`
	RoleName     string    `json:"role_name"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	IsActive     bool      `json:"is_active"`
	Permissions  []string  `json:"permissions"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	GetByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
	GetPermissionsByRoleID(roleID string) ([]string, error)
}
