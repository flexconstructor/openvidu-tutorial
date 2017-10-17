package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/flexconstructor/openvidu-tutorial/entity"
)

// mockOpenViDu is a mock that imitates the OpenViDu HTTP Client behavior.
type mockOpenViDu struct {
	behavior string
}

// GetMediaSession imitates HTTP Client GetMediaSession method behavior
// depending on one defined.
//
// Implements service.OpenViDu interface.
func (s *mockOpenViDu) GetMediaSession(sessionName string) (string, error) {
	if s.behavior == "ok" {
		return sessionName, nil
	}
	return "", errors.New("some error")
}

// GetToken imitates OpenViDu HTTP Client GetToken method behavior depending on
// one  defined.
//
// Implements service.OpenViDu interface.
func (s *mockOpenViDu) GetToken(
	params map[string]interface{}) (map[string]interface{}, error) {
	if s.behavior == "ok" {
		params["token"] = "test token"
		return params, nil
	}
	return nil, errors.New("some error")
}

// mockLoginAction is a mock that imitates the LoginAction behavior.
type mockLoginAction struct {
	behavior string
}

// Do imitates LoginAction Do method behavior depending on one
// defined.
func (a *mockLoginAction) Do(username string, password string) error {
	if a.behavior != "ok" {
		return errors.New("some error")
	}
	return nil
}

// mockSessionAction is a mock that imitates SessionAction behavior.
type mockSessionAction struct {
	behavior string
}

// Add imitates SessionAction Add method behavior depending on one
// defined.
func (a *mockSessionAction) Add(
	sessionID string, sessionName string, ownerName string) error {
	if a.behavior == "ok" {
		return nil
	}
	return errors.New("some error")
}

// Delete imitates SessionAction Delete method behavior depending on one
// defined.
func (a *mockSessionAction) Delete(sessionName string, userName string) error {
	if a.behavior == "ok" {
		return nil
	}
	return errors.New("some error")
}

// GetID imitates SessionAction GetID method behavior depending on one
// defined.
func (a *mockSessionAction) GetID(sessionName string) (string, error) {
	if a.behavior == "ok" {
		return "test session ID", nil
	}
	return "", errors.New("some error")
}

// IsExists imitates SessionAction IsExists method behavior depending on one
// defined.
func (a *mockSessionAction) IsExists(sessionName string) bool {

	return a.behavior == "ok"
}

func TestPages_Index(t *testing.T) {
	Convey("Writes index page to context", t, func() {
		_, ctx := newTestContext()
		(&Pages{}).Index(ctx)

		So(ctx.MustGet("template").(string), ShouldEqual, "index.tmpl")
	})

	Convey("Writes index page to context with error", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", nil)
		ctx.Error(errors.New("some error"))
		(&Pages{}).Index(ctx)

		So(ctx.MustGet("parameters").(gin.H), ShouldNotBeEmpty)
		So(ctx.MustGet("parameters").(gin.H)["error"], ShouldContainSubstring,
			"some error")
	})
}

func TestPages_Dashboard(t *testing.T) {
	Convey("Writes dashboard page to context", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("user", "test user")
		ctx.Request.PostForm.Add("pass", "test password")
		c := Pages{
			SessionStore:    &storeMock{behavior: "ok"},
			OpenViDuService: &mockOpenViDu{"ok"},
			LoginAction:     &mockLoginAction{"ok"},
		}
		session := sessions.NewSession(c.SessionStore, SESSION_NAME)
		session.Values = map[interface{}]interface{}{}
		c.SessionStore.Save(nil, nil, session)
		c.Dashboard(ctx)

		So(ctx.MustGet("template").(string), ShouldEqual, "dashboard.tmpl")
		So(ctx.Writer.Status(), ShouldEqual, http.StatusOK)

		Convey("writes user login to HTTP session", func() {
			So(session.Values["loggedUser"], ShouldEqual, "test user")
		})

		Convey("Session errors is nil", func() {
			So(session.Values["error"], ShouldBeNil)
		})
	})

	Convey("Redirect to index", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("user", "test user")
		ctx.Request.PostForm.Add("pass", "test password")
		(&Pages{SessionStore: &storeMock{behavior: "failure"},
			OpenViDuService: &mockOpenViDu{"ok"},
			LoginAction:     &mockLoginAction{"ok"},
		}).Dashboard(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusTemporaryRedirect)
	})

	Convey("If login action failed", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("user", "test user")
		ctx.Request.PostForm.Add("pass", "test password")
		c := Pages{
			SessionStore:    &storeMock{behavior: "ok"},
			OpenViDuService: &mockOpenViDu{"ok"},
			LoginAction:     &mockLoginAction{"failure"},
		}

		session := sessions.NewSession(c.SessionStore, SESSION_NAME)
		session.Values = map[interface{}]interface{}{}
		c.SessionStore.Save(nil, nil, session)
		c.Dashboard(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusTemporaryRedirect)

		Convey("Session error is not nil", func() {
			So(session.Values["error"], ShouldNotBeNil)
		})
	})
}

func TestPages_Session(t *testing.T) {
	Convey("Writes session data to context", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("session-name", "test session name")
		ctx.Request.PostForm.Add("data", "test session data")
		ctx.Set("user", &entity.User{Name: "test user name",
			Password: "test password", Role: 1})
		(&Pages{SessionAction: &mockSessionAction{"ok"},
			OpenViDuService: &mockOpenViDu{"ok"}}).Session(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusOK)

		Convey("Template is session.tmpl", func() {
			So(ctx.MustGet("template"), ShouldEqual, "session.tmpl")
		})

		Convey("Context error is nil", func() {
			So(ctx.Errors, ShouldBeEmpty)
		})

		Convey("With correct parameters", func() {
			params, ok := ctx.MustGet("parameters").(gin.H)
			So(ok, ShouldBeTrue)
			So(params["sessionId"], ShouldEqual, "test session ID")
			So(params["token"], ShouldEqual, "test token")
			So(params["userName"], ShouldEqual, "test user name")
			So(params["sessionName"], ShouldEqual, "test session name")
		})
	})

	Convey("If context has error", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Error(errors.New("some error"))
		(&Pages{}).Session(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusTemporaryRedirect)
	})

	Convey("if user can not publish", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("session-name", "test session name")
		ctx.Request.PostForm.Add("data", "test session data")
		ctx.Set("user", &entity.User{Name: "test user name",
			Password: "test password", Role: 0})
		(&Pages{SessionAction: &mockSessionAction{"failure"},
			OpenViDuService: &mockOpenViDu{"ok"}}).Session(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusTemporaryRedirect)
		So(ctx.Errors, ShouldNotBeEmpty)
	})

	Convey("If session action return error", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("session-name", "test session name")
		ctx.Request.PostForm.Add("data", "test session data")
		ctx.Set("user", &entity.User{Name: "test user name",
			Password: "test password", Role: 1})
		(&Pages{SessionAction: &mockSessionAction{"failure"},
			OpenViDuService: &mockOpenViDu{"ok"}}).Session(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusTemporaryRedirect)
		So(ctx.Errors, ShouldNotBeEmpty)
	})

	Convey("If openvidu service return an error", t, func() {
		_, ctx := newTestContext()
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("session-name", "test session name")
		ctx.Request.PostForm.Add("data", "test session data")
		ctx.Set("user", &entity.User{Name: "test user name",
			Password: "test password", Role: 1})
		(&Pages{SessionAction: &mockSessionAction{"ok"},
			OpenViDuService: &mockOpenViDu{"failure"}}).Session(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusTemporaryRedirect)
		So(ctx.Errors, ShouldNotBeEmpty)
	})
}

func TestPages_Leave(t *testing.T) {
	Convey("Leaves session", t, func() {
		_, ctx := newTestContext()
		ctx.Set("user", &entity.User{Name: "test user name",
			Password: "test password", Role: 1})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("session-name", "test session name")
		(&Pages{SessionAction: &mockSessionAction{"ok"}}).Leave(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusTemporaryRedirect)
		So(ctx.Errors, ShouldBeEmpty)
	})

	Convey("If leave session failed", t, func() {
		_, ctx := newTestContext()
		ctx.Set("user", &entity.User{Name: "test user name",
			Password: "test password", Role: 1})
		ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		ctx.Request.PostForm = url.Values{}
		ctx.Request.PostForm.Add("session-name", "test session name")
		(&Pages{SessionAction: &mockSessionAction{"failure"}}).Leave(ctx)

		So(ctx.Writer.Status(), ShouldEqual, http.StatusTemporaryRedirect)
		So(ctx.Errors, ShouldNotBeEmpty)
	})
}

// newTestContext initializes new HTTP request context and response recorder
// for test case.
func newTestContext() (w *httptest.ResponseRecorder, context *gin.Context) {
	w = httptest.NewRecorder()
	context, _ = gin.CreateTestContext(w)
	return
}
