package middlewares

import (
	"net/http"
	"simple-jwt-go/api/logger"
	"time"

	"github.com/sirupsen/logrus"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		lrw := logger.NewLog(w)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		entry := logrus.WithFields(logrus.Fields{
			"duration": duration,
			"method":   r.Method,
			"path":     r.RequestURI,
			"status":   lrw.StatusCode,
		})

		switch {
		case lrw.StatusCode > 0 && lrw.StatusCode < 400:
			entry.Info(string(lrw.Body))
		case lrw.StatusCode >= 400 && lrw.StatusCode < 500:
			entry.Warn(string(lrw.Body))
		case lrw.StatusCode >= 500:
			entry.Error(string(lrw.Body))
		default:
			entry.Warnf("unknown code: %d", lrw.StatusCode)
		}

	})
}
