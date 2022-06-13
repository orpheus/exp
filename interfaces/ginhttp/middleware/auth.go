package ginhttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	ginhttp "github.com/orpheus/exp/interfaces/ginhttp/auth"
	usecases "github.com/orpheus/exp/usecases/auth"
	"net/http"
	"strings"
)

func AuthGuardian(guardian ginhttp.PermissionGuardian) gin.HandlerFunc {
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

		token, err := usecases.JWTAuthService().ValidateToken(tokenString)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims, c.Request.RequestURI, c.Request.Method)
		requiredPermission := guardian.GetRequiredPermission(c.Request.RequestURI, c.Request.Method)
		if requiredPermission == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not find required permission"})
			return
		}

		scope := claims["scope"].([]interface{})
		hasPermission := false
		for _, p := range scope {
			hasPermission = ginhttp.HasPermission(requiredPermission, p.(string))
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("User does not have the following permission: %s", requiredPermission)})
			return
		}

		userId := claims["userId"]
		c.Set("userId", userId)

		c.Next()
	}
}
