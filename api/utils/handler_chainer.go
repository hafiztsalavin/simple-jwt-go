package utils

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func HandlerFuncs(middlewares []Middleware, handlerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	var handler http.Handler = http.HandlerFunc(handlerFunc)

	for i := len(middlewares); i > 0; i-- {
		handler = middlewares[i-1](handler)
	}

	return handler.ServeHTTP
}
