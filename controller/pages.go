package controller

import (
	"log"
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
}

// Index returns index page.
func (c *Pages) Index(ctx *gin.Context) {
	if ctx.Errors != nil && ctx.Request.Method == http.MethodPost {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{"error": ctx.Errors.String()})
	}
	ctx.HTML(http.StatusOK, "index.tmpl", nil)
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
	ctx.HTML(http.StatusOK, "dashboard.tmpl", nil)
}

// Session returns session page.
func (c *Pages) Session(ctx *gin.Context) {
	if ctx.Err() != nil {
		panic(ctx.Err())
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return
	}
	ovdSession := ctx.PostForm("session-name")
	participant := ctx.PostForm("data")
	log.Printf("Connect user: %s to session: %s", participant, ovdSession)
	session, err := c.OpenViDuService.GetMediaSession(ovdSession)
	if err != nil {
		ctx.Error(err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return

	}
	log.Printf("Session: %s", session)
	user := ctx.MustGet("user").(*entity.User)
	tokenOptions := make(map[string]interface{})
	tokenOptions["session"] = session
	tokenOptions["role"] = user.Role.String()
	tokenOptions["data"] = "DATA"
	tokenMap, err := c.OpenViDuService.GetToken(tokenOptions)
	if err != nil {
		ctx.Error(err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		ctx.Abort()
		return
	}

	ctx.HTML(http.StatusOK, "session.tmpl", gin.H{
		"sessionId":   tokenMap["session"],
		"token":       tokenMap["token"],
		"nickName":    participant,
		"userName":    user.Name,
		"sessionName": ovdSession,
	})
}
