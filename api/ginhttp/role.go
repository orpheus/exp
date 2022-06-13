package ginhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/system/sysauth"
	"net/http"
)

// RoleController Controller
type RoleController struct {
	Interactor RoleInteractor
	Logger     api.Logger
}

// RoleInteractor Interface. Lets the RoleController know what is can do
type RoleInteractor interface {
	FindAll() []sysauth.Role
	FindById(id uuid.UUID) (sysauth.Role, error)
	CreateOne(role sysauth.Role) (sysauth.Role, error)
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
	records := r.Interactor.FindAll()
	c.IndentedJSON(http.StatusOK, records)
}

func (r *RoleController) FindById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	record, err := r.Interactor.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, record)
}

func (r *RoleController) CreateOne(c *gin.Context) {
	var role sysauth.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec, err := r.Interactor.CreateOne(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
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
	err = r.Interactor.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, true)
}
