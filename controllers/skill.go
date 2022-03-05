package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/orpheus/exp/core"
	"github.com/orpheus/exp/repository"
	"log"
	"net/http"
	"strconv"
)

type SkillController struct {
	Router *gin.Engine
	Repo   *repository.SkillRepo
}

func (s *SkillController) RegisterRoutes() {
	router := s.Router.Group("/api")
	{
		router.GET("/skill", s.FindAllSkills)
		router.GET("/skill/:id", s.FindSkillById)
		router.POST("/skill", s.CreateSkill)
		router.POST("/skill/addTxp", s.AddTxp)
		router.DELETE("/skill/:id", s.DeleteById)
	}
}

func (s *SkillController) FindAllSkills(c *gin.Context) {
	Skills := s.Repo.FindAll()
	c.IndentedJSON(http.StatusOK, Skills)
}

func (s *SkillController) FindSkillById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", s, err)
	}
	log.Printf("successfully parsed UUID %v", id)
	skill, err := s.Repo.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, skill)
}

func (s *SkillController) CreateSkill(c *gin.Context) {
	var reqBody core.Skill
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

func (s *SkillController) AddTxp(c *gin.Context) {
	id := c.Query("id")
	txp := c.Query("txp")

	parsedId, err := uuid.FromString(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to parse id as uuid")
		return
	}

	parsedTxp, err := strconv.Atoi(txp)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to parse txp as integer")
		return
	}

	skill, err := s.Repo.FindById(parsedId)

	skill.AddTxp(parsedTxp)

	updatedSkill, err := s.Repo.UpdateExpLvl(skill)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}

	c.JSON(http.StatusOK, updatedSkill)
}

func (s *SkillController) DeleteById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", s, err)
	}
	response, err := s.Repo.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, response)
}
