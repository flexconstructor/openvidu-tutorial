package service

import "errors"

// OpenViDu is an interface of OpenViDu server.
type OpenViDu interface {
	// GetMediaSession calls OpenViDu server to retrieve of OpenViDu session
	// data object.
	GetMediaSession(sessionName string) (string, error)

	// GetToken calls OpenViDu server to retrieve of OpenViDu auth token
	// data object.
	GetToken(params map[string]interface{}) (map[string]interface{}, error)
}

// Service is an implementation of OpenViDu interface that performs retrieving
// of session and token from OpenViDu.
type Service struct {
	OpenViDu HTTPClient
}

//  GetMediaSession calls OpenViDu server to retrieve of OpenViDu session
// data object.
//
// Implements OpenViDu interface.
func (s *Service) GetMediaSession(
	sessionName string) (string, error) {
	m, err := s.OpenViDu.Post("api/sessions", nil)
	if err != nil {
		return "", err
	}
	session, ok := m["id"]
	if !ok {
		return "", errors.New("OpenViDu response contains no session ID")
	}
	sessionStr, ok := session.(string)

	if !ok {
		return "", errors.New("can not cast session ID to string")
	}
	return sessionStr, nil
}

// GetToken calls OpenViDu server to retrieve of OpenViDu auth token
// data object.
//
// Implements OpenViDu interface.
func (s *Service) GetToken(
	params map[string]interface{}) (map[string]interface{}, error) {
	m, err := s.OpenViDu.Post("api/tokens", params)
	if err != nil {
		return nil, err
	}
	return m, nil
}
