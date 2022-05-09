package routers

import (
	"simple-jwt-go/api/controllers"
	// "simple-jwt-go/api/middlewares"

	"github.com/gorilla/mux"
)

func AddUserRoutes(router *mux.Router) error {
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST").Name("SignUp")
	router.HandleFunc("/signin", controllers.SignIn).Methods("POST").Name("SignIn")

	return nil
}
