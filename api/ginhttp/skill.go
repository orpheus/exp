package ginhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/orpheus/exp/api"
	"github.com/orpheus/exp/core"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
)

// SkillController takes a service (SkillInteractor) and registers its
// method receivers as gin route handlers. It is fundamentally the job of this
// WebServiceHandler to handle incoming route requests and use the interactor
// service to handle the domain/business logic. The SkillController handlers
// will take care of request parsing.
type SkillController struct {
	Interactor SkillInteractor
	Logger     api.Logger
}

// SkillInteractor interface tells the SkillController what kind of object to expect.
// When our app is live, it will get constructed with the SkillInteractor struct from
// usecases which implements these methods.
type SkillInteractor interface {
	FindAllSkills() []core.Skill
	FindSkillById(id uuid.UUID) (core.Skill, error)
	CreateSkill(skill core.Skill, userId uuid.UUID) (core.Skill, error)
	DeleteById(skillId uuid.UUID) error
	ExistsBySkillIdAndUserId(skillId string, userId uuid.UUID) (bool, error)
	AddTxp(txp int, skillId uuid.UUID) (*core.Skill, error)
}

// RegisterRoutes takes a gin router group which determines the base
// path (i.e. /api).
func (s *SkillController) RegisterRoutes(router *gin.RouterGroup) {
	skill := router.Group("/skill")
	{
		skill.GET("", s.FindAllSkills)
		skill.GET("/:id", s.FindSkillById)
		skill.POST("", s.CreateSkill)
		skill.POST("/addTxp", s.AddTxp)
		skill.DELETE("/:id", s.DeleteById)
	}
}

func (s *SkillController) FindAllSkills(c *gin.Context) {
	skills := s.Interactor.FindAllSkills()
	c.JSON(http.StatusOK, skills)
}

func (s *SkillController) FindSkillById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		msg := fmt.Sprintf("failed to parse skill id. %s", err)
		c.JSON(http.StatusBadRequest, msg)
		return
	}
	s.Logger.Logf("successfully parsed UUID %v", id)

	skill, err := s.Interactor.FindSkillById(id)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, fmt.Sprintf("Skill %s not found", id))
			return
		}
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("%s", err.Error()))
		return
	}
	c.IndentedJSON(http.StatusOK, skill)
}

func (s *SkillController) CreateSkill(c *gin.Context) {
	var skill core.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
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

	rec, err := s.Interactor.CreateSkill(skill, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
		)
		return
	}

	c.JSON(http.StatusOK, rec)
}

func (s *SkillController) AddTxp(c *gin.Context) {
	queryId := c.Query("id")
	queryTxp := c.Query("txp")

	skillId, err := uuid.FromString(queryId)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to parse id as uuid")
		return
	}

	txp, err := strconv.Atoi(queryTxp)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Failed to parse txp as integer")
		return
	}

	updatedSkill, err := s.Interactor.AddTxp(txp, skillId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
		)
		return
	}

	c.JSON(http.StatusOK, updatedSkill)
}

func (s *SkillController) DeleteById(c *gin.Context) {
	skillId, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", s, err)
	}

	err = s.Interactor.DeleteById(skillId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, true)
}
