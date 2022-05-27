package routers

import (
	"simple-jwt-go/api/controllers"
	"simple-jwt-go/api/middlewares"
	"simple-jwt-go/api/utils"

	"github.com/gorilla/mux"
)

// Define routes
func Routes(router *mux.Router) error {
	router.HandleFunc("/register", utils.HandlerFuncs([]utils.Middleware{middlewares.Log}, controllers.SignUp)).Methods("POST").Name("SignUp")
	router.HandleFunc("/login", utils.HandlerFuncs([]utils.Middleware{middlewares.Log}, controllers.SignIn)).Methods("POST").Name("SignIn")

	return nil
}
