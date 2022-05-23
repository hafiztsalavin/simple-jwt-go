package middlewares

import (
	"net/http"
	"os"
	"simple-jwt-go/api/auth"
	"simple-jwt-go/api/models"
	"simple-jwt-go/api/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// check cookie is middleware that checks jwt token from incoming request has correct token or not
func CheckCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationCookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("authentication failed")), nil)
				return
			}
			utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
			return
		}

		jwtTokenString := authorizationCookie.Value
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
