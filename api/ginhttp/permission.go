package ginhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/system/sysauth"
	"net/http"
)

type PermissionController struct {
	Interactor PermissionInteractor
	Logger     api.Logger
}

type PermissionInteractor interface {
	FindAll() ([]sysauth.Permission, error)
	FindById(id string) (sysauth.Permission, error)
	CreateOne(id string) (sysauth.Permission, error)
	DeleteById(id string) error
	EnforcePermissions(pg sysauth.PermissionGetter)
}

func (p *PermissionController) RegisterRoutes(router *gin.RouterGroup) {
	permissions := router.Group("/permissions")
	{
		permissions.GET("", p.FindAll)
	}
}

func (p *PermissionController) FindAll(c *gin.Context) {
	permissions, err := p.Interactor.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}
