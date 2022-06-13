package ginhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/orpheus/exp/domain"
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
	service SkillInteractor
}

// SkillInteractor interface tells the SkillController what kind of object to expect.
// When our app is live, it will get constructed with the SkillInteractor struct from
// usecases which implements these methods.
type SkillInteractor interface {
	FindAllSkills() []domain.Skill
	FindSkillById(id uuid.UUID) (domain.Skill, error)
	CreateSkill(skill domain.Skill, userId uuid.UUID) (domain.Skill, error)
	DeleteById(skillId uuid.UUID) error
	ExistsByUserId(skillId uuid.UUID, userId uuid.UUID) (bool, error)
	AddTxp(txp int, skillId uuid.UUID) (*domain.Skill, error)
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
	skills := s.service.FindAllSkills()
	c.JSON(http.StatusOK, skills)
}

func (s *SkillController) FindSkillById(c *gin.Context) {
	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", s, err)
	}
	log.Printf("successfully parsed UUID %v", id)
	skill, err := s.service.FindSkillById(id)
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
	var skill domain.Skill
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

	skillUUID, _ := uuid.FromString(skill.SkillId)
	exists, err = s.service.ExistsByUserId(skillUUID, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Skill already exists for user"})
		return
	}

	rec, err := s.service.CreateSkill(skill, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
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

	updatedSkill, err := s.service.AddTxp(txp, skillId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
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

	err = s.service.DeleteById(skillId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("database error: %s", err.Error())},
		)
		return
	}
	c.JSON(http.StatusOK, true)
}
