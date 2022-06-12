package ginhttprouter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/interfaces/ginhttprouter/auth"
	"github.com/orpheus/exp/interfaces/persistence/repository"
	"net/http"
)

type PermissionController struct {
	Router   *gin.Engine
	Repo     *repository.PermissionRepo
	Guardian auth.PermissionGuardian
}

func (svc *PermissionController) RegisterRoutes(router *gin.RouterGroup) {
	permissions := router.Group("/permissions")
	{
		permissions.GET("/", svc.FindAll)
	}
}

func (svc *PermissionController) EnforcePermissions() {
	allPermissions := svc.Guardian.GetAllPermissions()
	existingPermissions, err := svc.Repo.FindAll()
	if err != nil {
		panic(err)
	}
	fmt.Println(existingPermissions)
	mappedPermissions := make(map[string]bool)
	for _, v := range existingPermissions {
		mappedPermissions[v.Id] = true
	}
	for _, p := range allPermissions {
		if !mappedPermissions[p] {
			_, err := svc.Repo.CreateOne(p)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (svc *PermissionController) FindAll(c *gin.Context) {
	permissions, err := svc.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}
