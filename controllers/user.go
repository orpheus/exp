package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/dto"
	"github.com/orpheus/exp/repository"
	"net/http"
)

type UserController struct {
	Router *gin.Engine
	Repo   *repository.UserRepo
}

func (svc *UserController) RegisterRoutes(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.GET("/", svc.FindAll)
	}
}

func (svc *UserController) FindAll(c *gin.Context) {
	users, err := svc.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userDTOs []dto.User
	for _, user := range users {
		userDTOs = append(userDTOs, user.ToDTO())
	}

	c.JSON(http.StatusOK, userDTOs)
}
