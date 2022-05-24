package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateJWTToken(userId uint32, username string) (string, error) {
	// 5 minute expired
	expirationTime := time.Now().Add(time.Minute * 1)
	claims := &Claims{
		ID:       userId,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_ACCESS_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(accessToken string) (string, error) {
	expirationTime := time.Now().Add(time.Minute * 5)
	claimsRefresh := &AccessClaims{
		AccessToken: accessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	refreshToken, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_KEY")))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
