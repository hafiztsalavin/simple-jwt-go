package auth

import (
	"os"
	"simple-jwt-go/api/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateJWTToken(userData models.User) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 30)
	claims := &Claims{
		ID:       userData.ID,
		Username: userData.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
