package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/flexconstructor/openvidu-tutorial/entity"
)

// storeMock  is a mock that imitates the HTTP Session storage behavior.
type storeMock struct {
	behavior        string
	sessionToReturn *sessions.Session
}

// Get imitates sessions.Store Get method behavior depending on one  defined.
//
// Implements sessions.Store interface.
func (s *storeMock) Get(
	r *http.Request, name string) (*sessions.Session, error) {
	if s.behavior == "failure" {
		return nil, errors.New("some error")
	}
	return s.sessionToReturn, nil
}

// New imitates sessions.Store New method behavior depending on one  defined.
//
// Implements sessions.Store interface.
func (s *storeMock) New(
	r *http.Request, name string) (*sessions.Session, error) {
	if s.behavior == "ok" {
		return sessions.NewSession(s, SESSION_NAME), nil
	}
	return nil, errors.New("some error")
}

// Save imitates sessions.Store Save method behavior depending on one  defined.
//
// Implements sessions.Store interface.
func (s *storeMock) Save(
	r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if s.behavior != "ok" {
		return errors.New("some error")
	}
	s.sessionToReturn = session
	return nil
}

// userActionMock is a mock that imitates UserAction behavior.
type userActionMock struct {
	behavior string
}

// Get imitates UserAction Get method behavior depending on one  defined.
func (a *userActionMock) Get(username string) (*entity.User, error) {
	if a.behavior == "ok" {
		return &entity.User{
			Name: "test user", Password: "test password", Role: 1}, nil
	}
	return nil, errors.New("some error")
}

func TestSession_Check(t *testing.T) {
	Convey("Writes user to context", t, func() {
		c := &Session{
			Store:      &storeMock{behavior: "ok"},
			UserAction: &userActionMock{behavior: "ok"},
		}
		c.Store.Save(nil, nil, &sessions.Session{
			Values: map[interface{}]interface{}{
				"loggedUser": "test user",
			},
		})
		_, ctx := runMiddlware(c.Check)

		So(func() { ctx.MustGet("user") }, ShouldNotPanic)

		Convey("with correct user data", func() {
			user := ctx.MustGet("user").(*entity.User)

			So(user.Name, ShouldEqual, "test user")
			So(user.Password, ShouldEqual, "test password")
			So(user.Role, ShouldEqual, 1)
		})

		Convey("Context errors should be empty", func() {
			So(ctx.Errors, ShouldBeEmpty)
		})
	})

	Convey("Writes session error to context", t, func() {
		c := &Session{
			Store:      &storeMock{behavior: "failure"},
			UserAction: &userActionMock{"ok"},
		}
		c.Store.Save(nil, nil, &sessions.Session{})
		_, ctx := runMiddlware(c.Check)

		So(ctx.Errors, ShouldNotBeEmpty)
	})

	Convey("If session contains some error", t, func() {
		c := &Session{
			Store:      &storeMock{behavior: "ok"},
			UserAction: &userActionMock{"ok"},
		}
		c.Store.Save(nil, nil, &sessions.Session{
			Values: map[interface{}]interface{}{
				"error": errors.New("some error"),
			},
		})
		_, ctx := runMiddlware(c.Check)

		So(ctx.Errors, ShouldNotBeEmpty)
	})

	Convey("If logged user not found", t, func() {
		c := &Session{
			Store:      &storeMock{behavior: "ok"},
			UserAction: &userActionMock{"ok"},
		}
		c.Store.Save(nil, nil, &sessions.Session{
			Values: map[interface{}]interface{}{},
		})
		_, ctx := runMiddlware(c.Check)

		So(ctx.Errors, ShouldNotBeEmpty)
	})

	Convey("If user not found", t, func() {
		c := &Session{
			Store:      &storeMock{behavior: "ok"},
			UserAction: &userActionMock{"failure"},
		}
		c.Store.Save(nil, nil, &sessions.Session{
			Values: map[interface{}]interface{}{
				"loggedUser": "test user",
			},
		})
		_, ctx := runMiddlware(c.Check)

		So(ctx.Errors, ShouldNotBeEmpty)
	})
}

// runMiddlware runs given middleware and returns its recorded response
// along with its final context.
func runMiddlware(fn gin.HandlerFunc) (
	w *httptest.ResponseRecorder,
	context *gin.Context,
) {
	w = httptest.NewRecorder()
	wait := make(chan struct{})
	g := gin.New()
	g.GET("/test", func(ctx *gin.Context) {
		defer close(wait)
		context = ctx
		ctx.Next() // wait to execute all other middlewares
	}, fn)
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	g.ServeHTTP(w, req)
	<-wait
	return
}
