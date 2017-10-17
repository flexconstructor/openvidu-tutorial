package service

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// httpClientMock is a mock that imitates the HTTP Client behavior.
type httpClientMock struct {
	behavior string
}

// Post imitates HTTP Client Post method behavior depending on one defined.
func (c *httpClientMock) Post(method string,
	args map[string]interface{}) (map[string]interface{}, error) {
	switch c.behavior {
	case "session":
		return map[string]interface{}{
			"id": "sessionID",
		}, nil
	case "no id":
		return map[string]interface{}{}, nil
	case "not string":
		return map[string]interface{}{"id": 1234}, nil
	case "token":
		return args, nil
	default:
		return nil, errors.New("some error")
	}

	return nil, nil
}

func TestService_GetMediaSession(t *testing.T) {
	Convey("Returns media session", t, func() {
		s := &Service{
			OpenViDu: &httpClientMock{"session"},
		}
		sessionID, err := s.GetMediaSession("session name")

		So(err, ShouldBeNil)
		So(sessionID, ShouldEqual, "sessionID")
	})

	Convey("Returns an error", t, func() {
		s := &Service{
			OpenViDu: &httpClientMock{"wrong"},
		}
		_, err := s.GetMediaSession("session name")

		So(err, ShouldNotBeNil)
	})

	Convey("if response contains no session id", t, func() {
		s := &Service{
			OpenViDu: &httpClientMock{"no id"},
		}
		_, err := s.GetMediaSession("session name")

		So(err, ShouldNotBeNil)
	})

	Convey("If session id is not string", t, func() {
		s := &Service{
			OpenViDu: &httpClientMock{"not string"},
		}
		_, err := s.GetMediaSession("session name")

		So(err, ShouldNotBeNil)
	})
}

func TestService_GetToken(t *testing.T) {
	Convey("Returns a map", t, func() {
		s := &Service{
			OpenViDu: &httpClientMock{"token"},
		}

		m, err := s.GetToken(map[string]interface{}{
			"test": "test",
		})

		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
		So(m, ShouldContainKey, "test")
	})

	Convey("Returns an error", t, func() {
		s := &Service{
			OpenViDu: &httpClientMock{"wrong"},
		}
		_, err := s.GetToken(nil)

		So(err, ShouldNotBeNil)
	})
}
