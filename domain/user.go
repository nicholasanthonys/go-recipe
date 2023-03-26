package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(context.Context, User) (User, error)
	Update(context.Context, UserID, User) (User, error)
	FindAll(context.Context) ([]User, error)
	FindByEmail(context.Context, Email) (User, error)
	FindByEmailAndPass(context context.Context, email string, password string) (User, error)
	FindByID(context.Context, UserID) (User, error)
	Delete(context.Context, UserID) error
}

type User struct {
	ID        UserID    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserID string
type Email string

func (r UserID) String() string {
	return string(r)
}

func (r Email) String() string {
	return string(r)
}

func NewUser(ID UserID, name string, password string, createdAt time.Time) User {
	return User{
		ID:        ID,
		Email:     name,
		Password:  password,
		CreatedAt: createdAt,
	}

}
