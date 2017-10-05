package entity

// UserRole is a role of user. Can be "SUBSCRIBER", "PUBLISHER" or "MODERATOR".
type UserRole uint8

// String defines string representation of user role.
func (r UserRole) String() string {
	roles := []string{"SUBSCRIBER", "PUBLISHER", "MODERATOR"}
	return roles[uint8(r)]
}

// User is a data of example`s user.
type User struct {
	Name     string
	Password string
	Role     UserRole
}

// Users is a repository interface that stores user data.
type Users interface {
	// Add adds users data to repository
	Add(username string, password string, role uint8)

	// Get retrieves user from repository.
	Get(username string) (*User, error)
}
