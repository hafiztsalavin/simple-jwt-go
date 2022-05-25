package middlewares

import (
	"net/http"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log
		next.ServeHTTP(w, r)
	})
}
