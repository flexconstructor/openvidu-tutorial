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
func InitRouter(HTTPClient service.HTTPClient) *gin.Engine {
	router := gin.Default()
	store := sessions.NewCookieStore([]byte("secret"))
	userRepo := repository.NewUsersRepository()
	userRepo.Add("publisher1", "pass", 1)
	userRepo.Add("publisher2", "pass", 1)
	userRepo.Add("subscriber", "pass", 0)

	s := &controller.Session{
		Store: store,
		UserAction: &action.User{
			UsersRepo: userRepo,
		},
	}
	router.Use(s.Check)
	router.Use(renderHTML)

	c := &controller.Pages{
		SessionStore: store,
		LoginAction: &action.Login{
			UserRepo: userRepo,
		},
		SessionAction: &action.Session{
			UserRepo:    userRepo,
			SessionRepo: repository.NewSessionsRepository(),
		},
		OpenViDuService: &service.Service{
			OpenViDu: HTTPClient,
		},
	}
	router.NoMethod(c.Index)
	router.NoRoute(c.Index)
	router.GET("/", c.Index)
	router.POST("/dashboard", c.Dashboard)
	router.POST("/session", c.Session)
	router.POST("/leave-session", c.Leave)
	return router
}
