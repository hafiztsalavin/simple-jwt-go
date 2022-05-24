package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"simple-jwt-go/api/auth"
	"simple-jwt-go/api/models"
	"simple-jwt-go/api/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// check cookie is middleware that checks jwt token from incoming request has correct token or not
func CheckAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationCookie, err := r.Cookie("access_token")
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
			return []byte(os.Getenv("JWT_ACCESS_KEY")), nil
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

func CheckRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		refreshCookie, err := r.Cookie("refresh_token")
		if err != nil {
			if err == http.ErrNoCookie {
				utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("authentication failed, please login")), nil)
				return
			}
			utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
			return
		}

		// get accessToken from Refresh Token
		jwtTokenString := refreshCookie.Value

		accessClaims := &auth.AccessClaims{}
		oldRefreshToken, err := jwt.ParseWithClaims(jwtTokenString, accessClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_REFRESH_KEY")), nil
		})

		fmt.Println("ini access token dr refresh: ", accessClaims.AccessToken)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("authentication failed, please login")), nil)
				return
			}
			utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse(err.Error())), nil)
			return
		}

		if !oldRefreshToken.Valid {
			utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("authentication failed, please login")), nil)
			return
		}

		// if valid, get payload id and username from old acces token
		if oldRefreshToken.Valid {

			// get Payload from Old Access Token
			claims := &auth.Claims{}
			jwtToken, err := jwt.ParseWithClaims(accessClaims.AccessToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_ACCESS_KEY")), nil
			})

			if jwtToken.Valid {
				utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("your token still valid")), nil)
				return
			}
			if err == nil {
				if err != jwt.ErrSignatureInvalid {
					utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("your token still valid")), nil)
					return
				}
				utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
				return
			}

			defer context.Clear(r)
			context.Set(r, "username", claims.Username)
			context.Set(r, "id", claims.ID)
		}
		next.ServeHTTP(w, r)
	})
}
