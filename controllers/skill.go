package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/orpheus/exp/core"
	"github.com/orpheus/exp/repository"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type SkillController struct {
	Router *gin.Engine
	Repo   *repository.SkillRepo
}

func (s *SkillController) RegisterRoutes(router *gin.RouterGroup) {
	skill := router.Group("/skill")
	{
		skill.GET("/", s.FindAllSkills)
		skill.GET("/:id", s.FindSkillById)
		skill.POST("/", s.CreateSkill)
		skill.POST("/addTxp", s.AddTxp)
		skill.DELETE("/:id", s.DeleteById)
	}
}

func (s *SkillController) FindAllSkills(c *gin.Context) {
	skills, err := s.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
		return
	}
	c.JSON(http.StatusOK, skills)
}

func (s *SkillController) FindSkillById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", s, err)
	}
	log.Printf("successfully parsed UUID %v", id)
	skill, err := s.Repo.FindById(id)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, fmt.Sprintf("Skill %s not found", id))
			return
		}
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

	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Could not find `userId` in auth token")},
		)
		return
	}

	userUuid, err := uuid.FromString(userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userId in token not valid uuid"})
		return
	}

	exists, err = s.Repo.ExistsByUserId(reqBody.SkillId, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Skill already exists for user"})
		return
	}

	reqBody.UserId = userUuid

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
	skillId := c.Query("id")
	txp := c.Query("txp")

	parsedSkillId, err := uuid.FromString(skillId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to parse id as uuid")
		return
	}

	parsedTxp, err := strconv.Atoi(txp)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to parse txp as integer")
		return
	}

	skill, err := s.Repo.FindById(parsedSkillId)

	if parsedTxp == 0 {
		c.JSON(http.StatusOK, skill)
		return
	}

	now := time.Now()
	last := skill.DateLastTxpAdd
	secondsSinceLastUpdate := int(now.Sub(last).Seconds())

	specialOneTimeBypass := specialRuleFirstTimeTxpApp(skill, parsedTxp)

	if parsedTxp > secondsSinceLastUpdate && !specialOneTimeBypass {
		var message string
		if skill.Txp == 0 {
			message = "Cannot add more than 3600 txp for the first hour of the skill's lifetime"
		} else {
			message = "Cannot add more txp than the difference of time in seconds between now and the last update"
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}

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

// Special Rule For New Skills: If the skill was created within the last hour and has 0 TXP, allow max 1 hour of txp
func specialRuleFirstTimeTxpApp(skill core.Skill, parsedTxp int) bool {
	now := time.Now()
	if skill.Txp == 0 && now.Sub(skill.DateCreated).Hours() < 1 {
		if parsedTxp <= int(time.Hour.Seconds()) {
			return true
		}
	}
	return false
}

func (s *SkillController) DeleteById(c *gin.Context) {
	skillId, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", s, err)
	}
	_, err = s.Repo.FindById(skillId)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, fmt.Sprintf("Skill %s not found", skillId))
			return
		}
	}
	response, err := s.Repo.DeleteById(skillId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, response)
}
