package ginhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/usecases/auth"
	"net/http"
)

// RoleController Controller
type RoleController struct {
	interactor RoleInteractor
}

// RoleInteractor Interface. Lets the RoleController know what is can do
type RoleInteractor interface {
	FindAll() []usecases.Role
	FindById(id uuid.UUID) (usecases.Role, error)
	CreateOne(role usecases.Role) (usecases.Role, error)
	DeleteById(id uuid.UUID) error
}

func (r *RoleController) RegisterRoutes(router *gin.RouterGroup) {
	role := router.Group("/role")
	{
		role.GET("/", r.FindAll)
		role.POST("/", r.CreateOne)
		role.GET("/:id", r.FindById)
		role.DELETE("/:id", r.DeleteById)
	}
}

func (r *RoleController) FindAll(c *gin.Context) {
	records := r.interactor.FindAll()
	c.IndentedJSON(http.StatusOK, records)
}

func (r *RoleController) FindById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	record, err := r.interactor.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, record)
}

func (r *RoleController) CreateOne(c *gin.Context) {
	var role usecases.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec, err := r.interactor.CreateOne(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, rec)
}

func (r *RoleController) DeleteById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = r.interactor.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, true)
}
