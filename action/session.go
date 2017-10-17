package action

import (
	"fmt"

	"github.com/flexconstructor/openvidu-tutorial/entity"
)

// Session is an action that performs operations with OpenViDu sessions.
type Session struct {
	SessionRepo entity.Sessions
	UserRepo    entity.Users
}

// Add adds new session with owner data but without participants.
//
// parameters:
//  sessionID   string   session ID that was returned from OpenViDu server.
//  sessionName string   The name of session that was returned from browser.
//  userName    string   Logged user name.
func (a *Session) Add(
	sessionID string, sessionName string, userName string) error {
	user, err := a.UserRepo.Get(userName)
	if err != nil {
		return err
	}

	if a.IsExists(sessionName) {
		_, err = a.addParticipant(sessionName, userName)
	} else {
		_, err = a.SessionRepo.Add(sessionID, sessionName, user)
	}
	return err
}

//  Delete delete participant of session i given userName is not name of
//  session`s owner or remove the session otherwise.
func (a *Session) Delete(sessionName string, userName string) error {
	user, err := a.UserRepo.Get(userName)
	if err != nil {
		return err
	}

	session, err := a.SessionRepo.Get(sessionName)
	if err != nil {
		return err
	}

	if session.Owner.Name == user.Name {
		return a.SessionRepo.Delete(sessionName)
	}

	return a.SessionRepo.Leave(sessionName, userName)
}

// GetID returns session ID by given session name.
func (a *Session) GetID(sessionName string) (string, error) {
	s, err := a.SessionRepo.Get(sessionName)
	if err != nil {
		return "", err
	}
	return s.ID, nil
}

// IsExists returns true if session is exists or false otherwise.
func (a *Session) IsExists(sessionName string) bool {
	_, err := a.SessionRepo.Get(sessionName)
	if err != nil {
		return false
	}
	return true
}

// addParticipant adds new participant to existed session.
func (a *Session) addParticipant(
	sessionName string, userName string) (string, error) {
	user, err := a.UserRepo.Get(userName)
	if err != nil {
		return "", err
	}

	session, err := a.SessionRepo.Get(sessionName)
	if err != nil {
		return "", err
	}

	if session.Owner.Name == userName {
		return "", fmt.Errorf(
			"owner %s can not subscribe session %s", userName, session.Name)
	}

	if _, ok := session.Subscribers[userName]; ok {
		return "", fmt.Errorf(
			"user %s already subscribed to the session %s", userName, sessionName)
	}
	session.AddParticipant(user)
	return session.ID, nil
}
