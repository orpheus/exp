package ginhttprouter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/domain"
	"net/http"
)

// SkillConfigController Controller.
type SkillConfigController struct {
	interactor SkillConfigInteractor
}

// SkillConfigInteractor Service Interface.
// TODO("Why do we define this here and not in the domain layer?")
type SkillConfigInteractor interface {
	FindAllSkillConfigs() []domain.SkillConfig
	FindSkillConfigById(id string) (domain.SkillConfig, error)
	CreateSkillConfig(skillConfig domain.SkillConfig) (domain.SkillConfig, error)
	CreateSkillConfigs(skillConfigs []domain.SkillConfig) error
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
	skillConfigs := s.interactor.FindAllSkillConfigs()
	c.IndentedJSON(http.StatusOK, skillConfigs)
}

func (s *SkillConfigController) FindSkillConfigById(c *gin.Context) {
	id := c.Param("id")
	skillConfigs, err := s.interactor.FindSkillConfigById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, skillConfigs)
}

func (s *SkillConfigController) CreateSkillConfig(c *gin.Context) {
	var reqBody domain.SkillConfig
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec, err := s.interactor.CreateSkillConfig(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, rec)
}

func (s *SkillConfigController) CreateSkillConfigs(c *gin.Context) {
	var reqBody []domain.SkillConfig
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.interactor.CreateSkillConfigs(reqBody)
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
	err := s.interactor.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, true)
}
