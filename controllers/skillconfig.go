package controllers

import (
	"com.orpheus/exp/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SkillConfigController struct {
	Router *gin.Engine
	Repo   *repository.SkillConfigRepo
}

func (s *SkillConfigController) RegisterRoutes() {
	router := s.Router.Group("/api")
	{
		router.GET("/skillConfig", s.FindAllSkillConfigs)
		router.GET("/skillConfig/:id", s.FindSkillConfigById)
		router.POST("/skillConfig", s.CreateSkillConfig)
		router.POST("/skillConfigs", s.CreateSkillConfigs)
		router.DELETE("/skillConfig/:id", s.DeleteById)
	}
}

func (s *SkillConfigController) FindAllSkillConfigs(c *gin.Context) {
	skillConfigs := s.Repo.FindAll()
	c.IndentedJSON(http.StatusOK, skillConfigs)
}

func (s *SkillConfigController) FindSkillConfigById(c *gin.Context) {
	id := c.Param("id")
	skillConfigs, err := s.Repo.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, skillConfigs)
}

func (s *SkillConfigController) CreateSkillConfig(c *gin.Context) {
	var reqBody repository.SkillConfig
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rec, err := s.Repo.CreateOne(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, rec)
}

func (s *SkillConfigController) CreateSkillConfigs(c *gin.Context) {
	var reqBody []repository.SkillConfig
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Req: ", reqBody)
	ret, err := s.Repo.CreateMany(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()},
		)
		return
	}
	c.JSON(http.StatusOK, ret)
}

func (s *SkillConfigController) DeleteById(c *gin.Context) {
	id := c.Param("id")
	response, err := s.Repo.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, response)
}
