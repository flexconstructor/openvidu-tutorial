package action

import "github.com/flexconstructor/openvidu-tutorial/entity"

// User is an action that returns saved user data from repository.
type User struct {
	UsersRepo entity.Users
}

// Get performs user data retrieving action.
func (a *User) Get(username string) (*entity.User, error) {
	user, err := a.UsersRepo.Get(username)
	if err != nil {
		return nil, err
	}
	return user, err
}
