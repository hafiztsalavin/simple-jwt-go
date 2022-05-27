package middlewares

import (
	"net/http"
	"simple-jwt-go/api/logger"
	"time"

	"github.com/sirupsen/logrus"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		// if err != nil {
		// 	fmt.Println("Could Not Open Log File : " + err.Error())
		// }

		start := time.Now()
		lrw := logger.NewLog(w)
		// logrus.SetFormatter(&logrus.JSONFormatter{})

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		entry := logrus.WithFields(logrus.Fields{
			"duration": duration,
			"method":   r.Method,
			"path":     r.RequestURI,
			"status":   lrw.StatusCode,
		})

		if lrw.StatusCode > 0 && lrw.StatusCode < 400 {
			entry.Info(string(lrw.Body))
		} else if lrw.StatusCode >= 400 && lrw.StatusCode < 500 {
			entry.Warn(string(lrw.Body))
		} else if lrw.StatusCode >= 500 {
			entry.Error(string(lrw.Body))
		} else {
			entry.Warnf("unknown code: %d", lrw.StatusCode)
		}

		// logrus.SetOutput(file)
	})
}
