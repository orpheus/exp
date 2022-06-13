package ginhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/system/sysauth"
	"net/http"
)

// SignOnController Controller
type SignOnController struct {
	Interactor SignOnInteractor
	Logger     api.Logger
}

type SignOnInteractor interface {
	Login(username string, password string) (sysauth.User, error)
	SignUp(user sysauth.User) (sysauth.User, error)
}

// RegisterRoutes registers a route group for login and signup apis
func (s *SignOnController) RegisterRoutes(router *gin.RouterGroup) {
	{
		router.POST("/login", s.Login)
		router.POST("/signup", s.SignUp)
	}
}

func (s *SignOnController) Login(c *gin.Context) {
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		c.JSON(http.StatusUnauthorized, "Missing required basic auth headers")
		return
	}

	user, err := s.Interactor.Login(username, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *SignOnController) SignUp(c *gin.Context) {
	var newUser sysauth.User

	err := c.ShouldBind(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.Interactor.SignUp(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
