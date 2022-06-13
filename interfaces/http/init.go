package ginhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/domain"
	"github.com/orpheus/exp/usecases"
	uca "github.com/orpheus/exp/usecases/auth"
)

func RegisterRoutes(r *gin.Engine) {
	//permissionGuardian := ginauth.MakePermissionGuardian()

	v1Router := r.Group("/api")

	//v1.Use(middleware.AuthGuardian(permissionGuardian))

	// V2 example
	// v2 := v1.Group("/v2")
	// v2.GET("/health", func() {}) /api/v2/health

	v1Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})

	// Role
	//roleRepo := &repository2.RoleRepo{DB: conn}
	//roleController := role.RoleController{
	//	Router: r,
	//	Repo:   roleRepo,
	//}
	//roleController.RegisterRoutes(v1)

	// Permissions
	permissionInteractor := uca.PermissionInteractor{}
	permissionController := PermissionController{
		interactor: &permissionInteractor,
	}
	permissionController.RegisterRoutes(v1Router)
	//permissionController.interactor.EnforcePermissions(permissionGuardian)

	// SignOn
	//signOnController := signon.SignOnController{
	//	Router:   r,
	//	Repo:     &repository2.UserRepo{DB: conn},
	//	RoleRepo: roleRepo,
	//	Auth:     auth2.JWTAuthService(),
	//}
	//signOnController.RegisterRoutes(v1)

	// SkillConfig
	skillConfigInteractor := usecases.SkillConfigInteractor{}
	skillConfigController := SkillConfigController{
		interactor: &skillConfigInteractor,
	}
	skillConfigController.RegisterRoutes(v1Router)

	// Skill
	skillInteractor := usecases.SkillInteractor{
		Policy: domain.SkillPolicyEnforcer{},
	}
	skillController := SkillController{
		service: &skillInteractor,
	}
	skillController.RegisterRoutes(v1Router)

	// User
	//userController := user.UserController{
	//	Router: r,
	//	Repo:   &repository2.UserRepo{DB: conn},
	//}
	//userController.RegisterRoutes(v1)
}