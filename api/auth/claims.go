package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID       uint32 `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
