package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orpheus/exp/auth"
	"github.com/orpheus/exp/middleware"
	"github.com/orpheus/exp/repository"
)

func RegisterAll(r *gin.Engine, conn *pgxpool.Pool) {
	permissionGuardian := auth.MakePermissionGuardian()

	v1 := r.Group("/api")
	v1.Use(middleware.AuthGuardian(permissionGuardian))

	// V2 example
	// v2 := v1.Group("/v2")
	// v2.GET("/health", func() {}) /api/v2/health

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})

	// Role
	roleRepo := &repository.RoleRepo{DB: conn}
	roleController := RoleController{
		Router: r,
		Repo:   roleRepo,
	}
	roleController.RegisterRoutes(v1)

	// Permissions
	permissionController := PermissionController{
		Router:   r,
		Repo:     &repository.PermissionRepo{DB: conn},
		Guardian: permissionGuardian,
	}
	permissionController.RegisterRoutes(v1)
	permissionController.EnforcePermissions()

	// SignOn
	signOnController := SignOnController{
		Router:   r,
		Repo:     &repository.UserRepo{DB: conn},
		RoleRepo: roleRepo,
		Auth:     auth.JWTAuthService(),
	}
	signOnController.RegisterRoutes(v1)

	// SkillConfig
	skillConfigController := SkillConfigController{
		Router: r,
		Repo:   &repository.SkillConfigRepo{DB: conn},
	}
	skillConfigController.RegisterRoutes(v1)

	// Skill
	skillController := SkillController{
		Router: r,
		Repo:   &repository.SkillRepo{DB: conn},
	}
	skillController.RegisterRoutes(v1)

	// User
	userController := UserController{
		Router: r,
		Repo:   &repository.UserRepo{DB: conn},
	}
	userController.RegisterRoutes(v1)
}
