package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"github.com/flexconstructor/openvidu-tutorial/entity"
)

const SESSION_NAME = "user_session"

// Session is a middleware that perform check in user data in HTTP session.
type Session struct {
	Store      sessions.Store
	UserAction interface {
		Get(username string) (*entity.User, error)
	}
}

// Check checks existed session and writes this to context.
func (mw *Session) Check(ctx *gin.Context) {
	s, err := mw.Store.Get(ctx.Request, SESSION_NAME)
	if err != nil {
		ctx.Error(err)
		return
	}

	e := s.Values["error"]
	if e != nil {
		ctx.Error(errors.New(e.(string)))
		return
	}

	username, ok := s.Values["loggedUser"].(string)

	if !ok {
		ctx.Error(errors.New("user not found"))
		return
	}

	user, err := mw.UserAction.Get(username)
	if err != nil {
		ctx.Error(errors.New("user not found"))
		return
	}
	ctx.Set("user", user)
}
