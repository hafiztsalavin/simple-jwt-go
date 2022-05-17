package middlewares

import (
	"net/http"
	"os"
	"simple-jwt-go/api/auth"
	"simple-jwt-go/api/models"
	"simple-jwt-go/api/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func CheckJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer ") {
			utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("authentication failed")), nil)
			return
		}

		jwtTokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		claims := &auth.Claims{}

		jwtToken, err := jwt.ParseWithClaims(jwtTokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("authentication failed")), nil)
				return
			}

			utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
			return
		}

		if !jwtToken.Valid {
			utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("authentication failed")), nil)
			return
		}

		defer context.Clear(r)
		context.Set(r, "username", claims.Username)
		context.Set(r, "id", claims.ID)

		next.ServeHTTP(w, r)
	})
}
