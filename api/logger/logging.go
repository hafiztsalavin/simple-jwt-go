package logger

import "net/http"

type ResponseData struct {
	Status int
}

type LoggingResponWriter struct {
	http.ResponseWriter
	ResponseData *ResponseData
}

func (r *LoggingResponWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	return size, err
}

func (r *LoggingResponWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode
}
