package repository

import (
	"fmt"

	"github.com/flexconstructor/openvidu-tutorial/entity"
)

// Sessions is a repository that stores OpenViDu sessions.
//
// implements entity.Sessions interface.
type Sessions struct {
	storage map[string]*entity.Session
}

// NewSessionsRepository returns new sessions repository instance.
//
// implements entity.Sessions interface.
func NewSessionsRepository() *Sessions {
	return &Sessions{
		storage: make(map[string]*entity.Session),
	}
}

// Add adds new session to repository.
//
// implements entity.Sessions interface.
func (r *Sessions) Add(
	sessionID string, sessionName string,
	owner *entity.User) (*entity.Session, error) {
	if _, ok := r.storage[sessionName]; ok {
		return nil, fmt.Errorf("session %s already exists", sessionName)
	}
	s := entity.NewSession()
	s.Name = sessionName
	s.ID = sessionID
	s.Owner = owner
	r.storage[sessionName] = s
	return s, nil
}

// Delete removes session from repository by given sessionName.
//
// implements entity.Sessions interface.
func (r *Sessions) Delete(sessionName string) error {
	if _, ok := r.storage[sessionName]; !ok {
		return fmt.Errorf("session %s does not exists", sessionName)
	}
	delete(r.storage, sessionName)
	return nil
}

// Get retrieves session from repository by given session name.
//
// implements entity.Sessions interface.
func (r *Sessions) Get(sessionName string) (*entity.Session, error) {
	s, ok := r.storage[sessionName]
	if !ok {
		return nil, fmt.Errorf("session %s does not exists", sessionName)
	}
	return s, nil
}

// Leave removes participant from session by given session name and user
// name.
//
// implements entity.Sessions interface.
func (r *Sessions) Leave(sessionName string, userName string) error {
	session, err := r.Get(sessionName)
	if err != nil {
		return err
	}
	if _, ok := session.Subscribers[userName]; !ok {
		return fmt.Errorf("user %s does not exists", userName)
	}
	delete(session.Subscribers, userName)
	return nil
}
