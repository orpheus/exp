package ginhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/usecases/auth"
	"net/http"
)

// UserController Controller
type UserController struct {
	Interactor UserInteractor
}

// UserInteractor Service Interface
type UserInteractor interface {
	FindAll() ([]*usecases.User, error)
}

// RegisterRoutes defines the route group for /user
func (u *UserController) RegisterRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("/", u.FindAll)
	}
}

// FindAll Returns users as defined by the UserInteractor implementation
func (u *UserController) FindAll(c *gin.Context) {
	users, err := u.Interactor.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, k := range users {
		k.RemovePassword()
	}

	c.JSON(http.StatusOK, users)
}
