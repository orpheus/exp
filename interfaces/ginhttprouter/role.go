package ginhttprouter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/interfaces/persistence/repository"
	"net/http"
)

type RoleController struct {
	Router *gin.Engine
	Repo   *repository.RoleRepo
}

func (svc *RoleController) RegisterRoutes(router *gin.RouterGroup) {
	role := router.Group("/role")
	{
		role.GET("/", svc.FindAll)
		role.POST("/", svc.CreateOne)
		role.GET("/:id", svc.FindById)
		role.DELETE("/:id", svc.DeleteById)
	}
}

func (svc *RoleController) FindAll(c *gin.Context) {
	records := svc.Repo.FindAll()
	c.IndentedJSON(http.StatusOK, records)
}

func (svc *RoleController) FindById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	record, err := svc.Repo.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, record)
}

func (svc *RoleController) CreateOne(c *gin.Context) {
	var reqBody repository.Role
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec, err := svc.Repo.CreateOne(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, rec)
}

func (svc *RoleController) DeleteById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response, err := svc.Repo.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, response)
}
