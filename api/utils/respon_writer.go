package utils

import (
	"encoding/json"
	"net/http"
)

func JSONResponseWriter(w *http.ResponseWriter, statusCode int, body interface{}, headerFields map[string]string) error {
	if body != nil {
		(*w).Header().Set("Content-Type", "application/json")
	}

	if headerFields != nil {
		responseHeader := (*w).Header()

		for k, v := range headerFields {
			responseHeader.Set(k, v)
		}
	}

	(*w).WriteHeader(statusCode)
	if body != nil {
		jsonString, err := json.Marshal(body)
		if err != nil {
			return err
		}

		(*w).Write(jsonString)
	}

	return nil
}
