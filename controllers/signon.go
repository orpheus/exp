package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orpheus/exp/auth"
	"github.com/orpheus/exp/repository"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type SignOnController struct {
	Router   *gin.Engine
	Repo     *repository.UserRepo
	RoleRepo *repository.RoleRepo
	Auth     auth.JWTService
}

func (svc *SignOnController) RegisterRoutes(router *gin.RouterGroup) {
	{
		router.POST("/login", svc.Login)
		router.POST("/signup", svc.SignUp)
	}
}

func (svc *SignOnController) Login(c *gin.Context) {
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		c.JSON(http.StatusUnauthorized, "Missing required basic auth headers")
		return
	}

	user, err := svc.Repo.GetUserLoginInfo(username)
	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// If the two passwords don't match, return a 401 status
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	role, err := svc.RoleRepo.FindById(user.RoleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Role %s does not exist", role.Id.String())})
		return
	}

	jwt := svc.Auth.GenerateToken(user.Id, role.Permissions)

	c.JSON(http.StatusOK, jwt)
}

func (svc *SignOnController) SignUp(c *gin.Context) {
	var newUser repository.User

	err := c.ShouldBind(&newUser) // create DTOs
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(newUser.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag"})
		return
	}

	_, err = svc.RoleRepo.FindById(newUser.RoleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Role %s does not exist", newUser.RoleId)})
		return
	}

	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	createdUser, err := svc.Repo.Create(newUser, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, createdUser.ToDTO())
}
