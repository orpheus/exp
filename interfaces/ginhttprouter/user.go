package ginhttprouter

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/usecases/auth"
	"net/http"
)

// UserController Controller
type UserController struct {
	interactor UserInteractor
}

// UserInteractor Service Interface
type UserInteractor interface {
	FindAll() ([]auth.User, error)
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
	users, err := u.interactor.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
