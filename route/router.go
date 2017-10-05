package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"github.com/flexconstructor/openvidu-tutorial/action"
	"github.com/flexconstructor/openvidu-tutorial/controller"
	"github.com/flexconstructor/openvidu-tutorial/repository"
	"github.com/flexconstructor/openvidu-tutorial/service"
)

// InitRouter initializes new HTTP router that performs routing of HTTP
// requests.
// Initializes all controllers.
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("/resources/templates/*.tmpl")
	router.Static("/images", "resources/static/images")
	router.StaticFile("/style.css", "resources/static//style.css")
	router.StaticFile("/openvidu-browser-1.1.0.js",
		"resources/static/openvidu-browser-1.1.0.js")
	store := sessions.NewCookieStore([]byte("secret"))
	repo := repository.NewUsersRepository()
	repo.Add("publisher1", "pass", 1)
	repo.Add("publisher2", "pass", 1)
	repo.Add("subscriber", "pass", 2)
	s := &controller.Session{
		Store: store,
		UserAction: &action.User{
			UsersRepo: repo,
		},
	}
	router.Use(s.Check)

	c := &controller.Pages{
		SessionStore: store,
		LoginAction: &action.Login{
			UserRepo: repo,
		},
		OpenViDuService: &service.Service{
			OpenViDu: &service.Client{
				OpenViDuURL: "https://openvidu-server-kms:8443",
				Login:       "OPENVIDUAPP",
				Password:    "MY_SECRET",
			},
		},
	}
	router.NoMethod(c.Index)
	router.NoRoute(c.Index)
	router.GET("/", c.Index)
	router.POST("/dashboard", c.Dashboard)
	router.POST("/session", c.Session)

	return router
}
