package entity

import (
	"fmt"
)

// Session is OpenViDu session value object performed by publisher for
// subscribers.
type Session struct {
	ID          string
	Name        string
	Owner       *User
	Subscribers map[string]*User
}

// NewSession returns new OpenViDu session value object.
func NewSession() *Session {
	return &Session{
		Subscribers: make(map[string]*User),
	}
}

// AddParticipant adds participant to session subscribers list.
func (e *Session) AddParticipant(user *User) {
	e.Subscribers[user.Name] = user
}

// RemoveParticipant removes participant from session subscribers list.
func (e *Session) RemoveParticipant(user *User) error {
	if _, ok := e.Subscribers[user.Name]; !ok {
		return fmt.Errorf("subscriber: %s not found", user.Name)
	}
	delete(e.Subscribers, user.Name)
	return nil
}

// Sessions is a repository that stores OpenViDu sessions.
type Sessions interface {

	// Add new session to repository by given session ID, session name, and
	// owner value object.
	Add(sessionID string, sessionName string, owner *User) (*Session, error)

	// Delete session by given session name.
	Delete(sessionName string) error

	// Get returns session by given session name.
	Get(sessionName string) (*Session, error)

	// Leave removes participant from session by given session name and user
	// name.
	Leave(sessionName string, userName string) error
}
