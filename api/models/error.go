package models

type ErrorResponse struct {
	Message string `json:"message" example:"failure to connect to db"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{Message: message}
}
