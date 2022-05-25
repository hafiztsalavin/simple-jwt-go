package auth

import (
	"github.com/dgrijalva/jwt-go"
)

// result as an auth token
type Claims struct {
	Id    uint32
	Email string
	jwt.StandardClaims
}

type AccessClaims struct {
	AccessToken string
	jwt.StandardClaims
}
