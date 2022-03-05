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
	Router *gin.Engine
	Repo   *repository.UserRepo
	Auth   auth.JWTService
}

func (svc *SignOnController) RegisterRoutes() {
	router := svc.Router.Group("/api")
	{
		router.POST("/login", svc.Login)
		router.POST("/signup", svc.SignUp)
	}
}

type LoginCredentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (svc *SignOnController) Login(c *gin.Context) {
	var credentials LoginCredentials
	err := c.ShouldBind(&credentials)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password, err := svc.Repo.GetPasswordForUser(credentials.Username)
	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(credentials.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	jwt := svc.Auth.GenerateToken(credentials.Username, true)

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

	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	fmt.Println(hashedPassword)
	createdUser, err := svc.Repo.Create(newUser, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, createdUser)
}