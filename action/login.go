package action

import (
	"errors"

	"github.com/flexconstructor/openvidu-tutorial/entity"
)

// Login is an action that performs authorization of user with login and
// password.
type Login struct {
	UserRepo entity.Users
}

// Do performs authorization action for given login and password.
// To authorize user password must be correct.
func (a *Login) Do(username string, password string) error {
	user, err := a.UserRepo.Get(username)
	if err != nil {
		return err
	}
	if user.Password != password {
		return errors.New("password incorrect")
	}
	return nil
}
