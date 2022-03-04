package controllers

import (
	"com.orpheus/exp/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RegisterAll(r *gin.Engine, conn *pgxpool.Pool) {
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})

	skillConfigController := SkillConfigController{
		Router: r,
		Repo:   &repository.SkillConfigRepo{DB: conn},
	}
	skillConfigController.RegisterRoutes()
}
