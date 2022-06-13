package ginhttp

import (
	"github.com/gin-gonic/gin"
	ucauth "github.com/orpheus/exp/usecases/auth"
	"net/http"
)

type PermissionController struct {
	interactor PermissionInteractor
}

type PermissionInteractor interface {
	FindAll() ([]ucauth.Permission, error)
	FindById(id string) (ucauth.Permission, error)
	CreateOne(id string) (ucauth.Permission, error)
	DeleteById(id string) error
	EnforcePermissions(pg ucauth.PermissionGetter)
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