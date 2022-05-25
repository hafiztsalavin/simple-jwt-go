package middlewares

import (
	"log"
	"net/http"
	"time"

	"simple-jwt-go/api/logger"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := logger.NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		statusCode := lrw.StatusCode
		log.Printf(
			"%s\t%s\t%d\t%s",
			r.Method,
			r.RequestURI,
			statusCode,
			time.Since(start),
		)
	})
}
