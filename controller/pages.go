package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"github.com/flexconstructor/openvidu-tutorial/entity"
	"github.com/flexconstructor/openvidu-tutorial/service"
)

// Pages is a HTTP controller that provides operations with HTTP pages of the
// example.
type Pages struct {
	SessionStore    sessions.Store
	OpenViDuService service.OpenViDu
	LoginAction     interface {
		Do(username string, password string) error
	}

	SessionAction interface {
		Add(sessionID string, sessionName string, ownerName string) error
		Delete(sessionName string, userName string) error
		GetID(sessionName string) (string, error)
		IsExists(sessionName string) bool
	}
}

// Index returns index page.
func (c *Pages) Index(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
	ctx.Set("template", "index.tmpl")
	if ctx.Errors != nil && ctx.Request.Method == http.MethodPost {
		ctx.Set("parameters", gin.H{"error": ctx.Errors.String()})
		return
	}
	ctx.Set("parameters", gin.H{})
}

//Dashboard returns dashboard page.
func (c *Pages) Dashboard(ctx *gin.Context) {
	session, err := c.SessionStore.Get(ctx.Request, SESSION_NAME)
	if err != nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return
	}

	login := ctx.PostForm("user")
	password := ctx.PostForm("pass")

	err = c.LoginAction.Do(login, password)
	if err != nil {
		session.Values["error"] = err.Error()
		session.Save(ctx.Request, ctx.Writer)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	session.Values["error"] = nil
	session.Values["loggedUser"] = login
	session.Save(ctx.Request, ctx.Writer)

	ctx.Status(http.StatusOK)
	ctx.Set("template", "dashboard.tmpl")
	ctx.Set("parameters", gin.H{})
}

// Session returns session page.
func (c *Pages) Session(ctx *gin.Context) {
	if len(ctx.Errors) > 0 {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return
	}
	sessionName := ctx.PostForm("session-name")
	participant := ctx.PostForm("data")
	user := ctx.MustGet("user").(*entity.User)
	var session string
	var err error
	if c.SessionAction.IsExists(sessionName) {
		session, err = c.SessionAction.GetID(sessionName)
	} else if user.Role > 0 {
		session, err = c.OpenViDuService.GetMediaSession(sessionName)
	} else {
		err = fmt.Errorf("user %s can not publish", participant)
	}

	if err != nil {
		ctx.Error(err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return
	}

	tokenOptions := make(map[string]interface{})
	tokenOptions["session"] = session
	tokenOptions["role"] = user.Role.String()
	tokenOptions["data"] = "{\"serverData\": \"" + participant + "\"}"
	tokenMap, err := c.OpenViDuService.GetToken(tokenOptions)
	if err != nil {
		ctx.Error(err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return
	}

	err = c.SessionAction.Add(session, sessionName, user.Name)
	if err != nil {
		ctx.Error(err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return
	}
	ctx.Status(http.StatusOK)
	ctx.Set("template", "session.tmpl")
	ctx.Set("parameters", gin.H{
		"sessionId":   tokenMap["session"],
		"token":       tokenMap["token"],
		"nickName":    participant,
		"userName":    user.Name,
		"sessionName": sessionName,
	})
}

// Leave the controller command that removes user from the OpenViDu session, or
// removes session if user is owner.
func (c *Pages) Leave(ctx *gin.Context) {
	ovdSession := ctx.PostForm("session-name")
	user := ctx.MustGet("user").(*entity.User)
	err := c.SessionAction.Delete(ovdSession, user.Name)
	if err != nil {
		ctx.Error(err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}
