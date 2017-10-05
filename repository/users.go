package repository

import (
	"errors"

	"github.com/flexconstructor/openvidu-tutorial/entity"
)

// Users implementation of entity.Users repository.
type Users struct {
	users map[string]*entity.User
}

// NewUsersRepository returns new instance of Users repository.
func NewUsersRepository() *Users {
	return &Users{
		users: make(map[string]*entity.User),
	}
}

// Add adds user data to repository.
//
// Implements entity.Users interface.
func (r *Users) Add(username string, password string, role uint8) {
	r.users[username] = &entity.User{
		Name:     username,
		Password: password,
		Role:     entity.UserRole(role),
	}
}

// Get retrieves user from repository.
//
// Implements entity.Users interface.
func (r *Users) Get(username string) (*entity.User, error) {
	user, ok := r.users[username]
	switch {
	case !ok:
		return nil, errors.New("Login incorrect")
	default:
		return user, nil
	}
}
