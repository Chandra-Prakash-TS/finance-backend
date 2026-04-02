package domain

import (
	"context"
	"time"
)

const (
	RoleViewer  = "viewer"
	RoleAnalyst = "analyst"
	RoleAdmin   = "admin"
)

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, page, pageSize int) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, id string) error
}
