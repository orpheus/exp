package auth

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

type JWTService interface {
	GenerateToken(userId uuid.UUID, scope []string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type authCustomClaims struct {
	UserId string   `json:"userId"`
	Scope  []string `json:"scope"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issuer    string
}

func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issuer:    "github.com/orpheus/exp",
	}
}

func getSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (service *jwtServices) GenerateToken(userId uuid.UUID, scope []string) string {
	claims := &authCustomClaims{
		userId.String(),
		scope,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token %s", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
}
