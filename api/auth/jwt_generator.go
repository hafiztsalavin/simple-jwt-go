package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateJWTToken(userId uint32, email string, exp time.Time) (string, error) {
	claims := &Claims{
		Id:    userId,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(accessToken string, exp time.Time) (string, error) {
	claimsRefresh := &AccessClaims{
		AccessToken: accessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	refreshToken, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_KEY")))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
