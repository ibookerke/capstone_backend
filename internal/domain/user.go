package domain

import (
	"context"
	"time"
)

type UserID int64

type User struct {
	ID                   UserID                 `json:"id,omitempty"`
	Email                string                 `json:"email,omitempty"`
	Password             string                 `json:"-"`
	Name                 string                 `json:"name,omitempty"`
	Surname              string                 `json:"surname,omitempty"`
	CreatedAt            time.Time              `json:"createdAt,omitempty"`
	UpdatedAt            time.Time              `json:"updatedAt,omitempty"`
	NotificationSettings []NotificationSettings `json:"notificationSettings,omitempty"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id UserID) (*User, error)
}
