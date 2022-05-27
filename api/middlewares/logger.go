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

		responseData := &logger.ResponseData{
			Status: 0,
		}

		lrw := logger.LoggingResponWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		next.ServeHTTP(&lrw, r)

		duration := time.Since(start)

		entry := logrus.WithFields(logrus.Fields{
			"duration": duration,
			"method":   r.Method,
			"path":     r.RequestURI,
			"status":   responseData.Status,
		})

		switch {
		case lrw.ResponseData.Status > 0 && lrw.ResponseData.Status < 400:
			entry.Info()
		case lrw.ResponseData.Status >= 400 && lrw.ResponseData.Status < 500:
			entry.Warn()
		case lrw.ResponseData.Status >= 500:
			entry.Error()
		default:
			entry.Warnf("unknown code: %d", lrw.ResponseData.Status)
		}

	})
}
