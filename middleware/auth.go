package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/orpheus/exp/auth"
	"net/http"
	"strings"
)

func AuthGuardian() gin.HandlerFunc {
	PermissionGuardian := auth.MakePermissionGuardian()
	return func(c *gin.Context) {
		if PermissionGuardian.HasOpenPermission(c.Request.RequestURI, c.Request.Method) {
			return
		}

		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		const BearerSchema = "Bearer"
		if !strings.Contains(authHeader, BearerSchema) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing `Bearer` in Authorization header"})
			return
		}
		tokenString := strings.TrimSpace(authHeader[len(BearerSchema):])
		fmt.Println(tokenString)

		token, err := auth.JWTAuthService().ValidateToken(tokenString)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(token)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println(claims)
		}

	}
}
