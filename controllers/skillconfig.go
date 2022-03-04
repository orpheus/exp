package controllers

import (
	"com.orpheus/exp/repository"
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
	}
}

func (s *SkillConfigController) FindAllSkillConfigs(c *gin.Context) {
	skillConfigs := s.Repo.FindAll()
	c.IndentedJSON(http.StatusOK, skillConfigs)
}
