package ginhttprouter

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/interfaces/ginhttprouter/auth"
	auth2 "github.com/orpheus/exp/usecases/auth"
	"net/http"
)

type PermissionController struct {
	interactor PermissionInteractor
	guardian   auth.PermissionGuardian
}

type PermissionInteractor interface {
	FindAll() ([]auth2.Permission, error)
	FindById(id string) (auth2.Permission, error)
	CreateOne(id string) (auth2.Permission, error)
	DeleteById(id string) error
}

func (p *PermissionController) RegisterRoutes(router *gin.RouterGroup) {
	permissions := router.Group("/permissions")
	{
		permissions.GET("/", p.FindAll)
	}
}

func (p *PermissionController) FindAll(c *gin.Context) {
	permissions, err := p.interactor.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}
