package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orpheus/exp/domain"
	"github.com/orpheus/exp/interfaces/ginhttp"
	"github.com/orpheus/exp/interfaces/ginhttp/auth"
	"github.com/orpheus/exp/interfaces/ginhttp/middleware"
	"github.com/orpheus/exp/interfaces/repository"
	"github.com/orpheus/exp/usecases"
	authuc "github.com/orpheus/exp/usecases/auth"
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
	jwtService := authuc.JWTAuthService()

	v1Router := r.Group("/api")
	v1Router.Use(middleware.AuthGuardian(permissionGuardian))
	v1Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "healthy",
		})
	})

	// Repositories
	roleRepository := &repository.RoleRepository{
		DB:     conn,
		Logger: tmpLogger,
	}
	permissionRepository := &repository.PermissionRepository{
		DB:     conn,
		Logger: tmpLogger,
	}
	skillConfigRepository := &repository.SkillConfigRepository{
		DB:     conn,
		Logger: tmpLogger,
	}
	skillRepository := &repository.SkillRepository{
		DB:     conn,
		Logger: tmpLogger,
	}
	skillerRepository := &repository.SkillerRepository{
		DB:     conn,
		Logger: tmpLogger,
	}

	// Interactors
	roleInteractor := &authuc.RoleInteractor{
		RoleRepository: roleRepository,
		Logger:         tmpLogger,
	}

	permissionInteractor := &authuc.PermissionInteractor{
		PermissionRepository: permissionRepository,
		Logger:               tmpLogger,
	}

	skillConfigInteractor := &usecases.SkillConfigInteractor{
		Repo:   skillConfigRepository,
		Logger: tmpLogger,
	}

	skillInteractor := &usecases.SkillInteractor{
		SkillRepository: skillRepository,
		Policy:          &domain.SkillPolicyEnforcer{},
		Logger:          tmpLogger,
	}

	skillerInteractor := &usecases.SkillerInteractor{
		Repo:   skillerRepository,
		Logger: tmpLogger,
	}

	userInteractor := &authuc.UserInteractor{
		SkillerInteractor: skillerInteractor,
		Logger:            tmpLogger,
	}

	signonInteractor := &authuc.SignOnInteractor{
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
