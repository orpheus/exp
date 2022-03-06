package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/orpheus/exp/auth"
	"net/http"
	"strings"
)

func AuthGuardian(guardian auth.PermissionGuardian) gin.HandlerFunc {
	return func(c *gin.Context) {
		if guardian.HasOpenPermission(c.Request.RequestURI, c.Request.Method) {
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

		token, err := auth.JWTAuthService().ValidateToken(tokenString)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		requiredPermission := guardian.GetRequiredPermission(c.Request.RequestURI, c.Request.Method)
		if requiredPermission == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not find required permission"})
			return
		}
		scope := claims["scope"].([]interface{})
		hasPermission := false
		for _, p := range scope {
			if requiredPermission == p {
				hasPermission = true
				break
			}
			if p == "*" {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("User does not have the following permission: %s", requiredPermission)})
			return
		}

	}
}