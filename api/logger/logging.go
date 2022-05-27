package logger

import "net/http"

type LoggingResponWriter struct {
	wrapped    http.ResponseWriter
	StatusCode int
	Body       []byte
}

func NewLog(wrapped http.ResponseWriter) *LoggingResponWriter {
	return &LoggingResponWriter{wrapped: wrapped}
}

func (r *LoggingResponWriter) Header() http.Header {
	return r.wrapped.Header()
}

func (r *LoggingResponWriter) Write(b []byte) (int, error) {
	if r.StatusCode >= 400 {
		r.Body = append(r.Body, b...)
	}

	return r.wrapped.Write(b)
}

func (r *LoggingResponWriter) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.wrapped.WriteHeader(statusCode)
}
