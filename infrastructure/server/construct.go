package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orpheus/exp/api/ginhttp"
	"github.com/orpheus/exp/api/ginhttp/auth"
	"github.com/orpheus/exp/api/ginhttp/middleware"
	"github.com/orpheus/exp/api/pgxrepo"
	"github.com/orpheus/exp/core"
	"github.com/orpheus/exp/system/sysauth"

	"github.com/orpheus/exp/system"

	"log"
)

type TmpLogger struct{}

func (l *TmpLogger) Log(v ...interface{}) {
	log.Println(v)
}

func (l *TmpLogger) Logf(format string, v ...interface{}) {
	log.Printf(format, v)
}

func Construct(r *gin.Engine, conn *pgxpool.Pool) {
	tmpLogger := &TmpLogger{}

	permissionGuardian := auth.MakePermissionGuardian()
	jwtService := sysauth.JWTAuthService()

	v1Router := r.Group("/api")
	v1Router.Use(middleware.AuthGuardian(permissionGuardian))
	v1Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})

	// Repositories
	roleRepository := &pgxrepo.RoleRepository{
		DB:     conn,
		Logger: tmpLogger,
	}
	permissionRepository := &pgxrepo.PermissionRepository{
		DB:     conn,
		Logger: tmpLogger,
	}
	skillConfigRepository := &pgxrepo.SkillConfigRepository{
		DB:     conn,
		Logger: tmpLogger,
	}
	skillRepository := &pgxrepo.SkillRepository{
		DB:     conn,
		Logger: tmpLogger,
	}
	skillerRepository := &pgxrepo.SkillerRepository{
		DB:     conn,
		Logger: tmpLogger,
	}

	// Interactors
	roleInteractor := &sysauth.RoleInteractor{
		RoleRepository: roleRepository,
		Logger:         tmpLogger,
	}

	permissionInteractor := &sysauth.PermissionInteractor{
		PermissionRepository: permissionRepository,
		Logger:               tmpLogger,
	}

	skillConfigInteractor := &system.SkillConfigInteractor{
		Repo:   skillConfigRepository,
		Logger: tmpLogger,
	}

	skillInteractor := &system.SkillInteractor{
		SkillRepository: skillRepository,
		Policy:          &core.SkillPolicyEnforcer{},
		Logger:          tmpLogger,
	}

	skillerInteractor := &system.SkillerInteractor{
		Repo:   skillerRepository,
		Logger: tmpLogger,
	}

	userInteractor := &sysauth.UserInteractor{
		SkillerInteractor: skillerInteractor,
		Logger:            tmpLogger,
	}

	signonInteractor := &sysauth.SignOnInteractor{
		UserRepository: userInteractor,
		RoleRepository: roleRepository,
		JWTService:     jwtService,
		Logger:         tmpLogger,
	}

	// Controllers
	permissionController := ginhttp.PermissionController{
		Interactor: permissionInteractor,
		Logger:     tmpLogger,
	}
	permissionController.RegisterRoutes(v1Router)
	permissionController.Interactor.EnforcePermissions(permissionGuardian)

	roleController := ginhttp.RoleController{
		Interactor: roleInteractor,
		Logger:     tmpLogger,
	}
	roleController.RegisterRoutes(v1Router)

	signOnController := ginhttp.SignOnController{
		Interactor: signonInteractor,
		Logger:     tmpLogger,
	}
	signOnController.RegisterRoutes(v1Router)

	skillConfigController := ginhttp.SkillConfigController{
		Interactor: skillConfigInteractor,
		Logger:     tmpLogger,
	}
	skillConfigController.RegisterRoutes(v1Router)

	skillController := ginhttp.SkillController{
		Interactor: skillInteractor,
		Logger:     tmpLogger,
	}
	skillController.RegisterRoutes(v1Router)

	userController := ginhttp.UserController{
		Interactor: userInteractor,
	}
	userController.RegisterRoutes(v1Router)
}
