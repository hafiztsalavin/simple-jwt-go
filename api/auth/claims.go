package auth

import (
	"github.com/dgrijalva/jwt-go"
)

// result as an auth token
type Claims struct {
	ID       uint32 `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
