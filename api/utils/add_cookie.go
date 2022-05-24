package utils

import (
	"net/http"
	"time"
)

func AddCookie(response http.ResponseWriter, name, value string, exp time.Time) {
	refreshCookie := &http.Cookie{
		Name:    name,
		Value:   value,
		Expires: exp,
	}
	http.SetCookie(response, refreshCookie)
}
