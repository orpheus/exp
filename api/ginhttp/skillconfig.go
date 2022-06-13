package ginhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"net/http"
)

// SkillConfigController Controller.
type SkillConfigController struct {
	Interactor SkillConfigInteractor
	Logger     api.Logger
}

// SkillConfigInteractor Service Interface.
// TODO("Why do we define this here and not in the domain layer?")
type SkillConfigInteractor interface {
	FindAllSkillConfigs() []core.SkillConfig
	FindSkillConfigById(id string) (core.SkillConfig, error)
	CreateSkillConfig(skillConfig core.SkillConfig) (core.SkillConfig, error)
	CreateSkillConfigs(skillConfigs []core.SkillConfig) error
	DeleteById(id string) error
}

func (s *SkillConfigController) RegisterRoutes(router *gin.RouterGroup) {
	skillConfig := router.Group("/skillConfig")
	{
		skillConfig.GET("/", s.FindAllSkillConfigs)
		skillConfig.POST("/", s.CreateSkillConfig)
		skillConfig.GET("/:id", s.FindSkillConfigById)
		skillConfig.DELETE("/:id", s.DeleteById)
	}

	skillConfigs := router.Group("skillConfigs")
	{
		skillConfigs.POST("/", s.CreateSkillConfigs)
	}
}

func (s *SkillConfigController) FindAllSkillConfigs(c *gin.Context) {
	skillConfigs := s.Interactor.FindAllSkillConfigs()
	c.IndentedJSON(http.StatusOK, skillConfigs)
}

func (s *SkillConfigController) FindSkillConfigById(c *gin.Context) {
	id := c.Param("id")
	skillConfigs, err := s.Interactor.FindSkillConfigById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, skillConfigs)
}

func (s *SkillConfigController) CreateSkillConfig(c *gin.Context) {
	var reqBody core.SkillConfig
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec, err := s.Interactor.CreateSkillConfig(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, rec)
}

func (s *SkillConfigController) CreateSkillConfigs(c *gin.Context) {
	var reqBody []core.SkillConfig
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.Interactor.CreateSkillConfigs(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, "Created")
}

func (s *SkillConfigController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	err := s.Interactor.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, true)
}
