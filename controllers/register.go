package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orpheus/exp/auth"
	"github.com/orpheus/exp/repository"
)

func RegisterAll(r *gin.Engine, conn *pgxpool.Pool) {
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})

	// SkillConfig
	skillConfigController := SkillConfigController{
		Router: r,
		Repo:   &repository.SkillConfigRepo{DB: conn},
	}
	skillConfigController.RegisterRoutes()

	// Skill
	skillController := SkillController{
		Router: r,
		Repo:   &repository.SkillRepo{DB: conn},
	}
	skillController.RegisterRoutes()

	// SignOn
	signOnController := SignOnController{
		Router: r,
		Repo:   &repository.UserRepo{DB: conn},
		Auth:   auth.JWTAuthService(),
	}
	signOnController.RegisterRoutes()
}
