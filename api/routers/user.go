package routers

import (
	"simple-jwt-go/api/controllers"
	// "simple-jwt-go/api/middlewares"

	"github.com/gorilla/mux"
)

func AddUserRoutes(router *mux.Router) error {
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST").Name("SignUp")

	// router.HandleFunc("/signup",utils.ChainHandlerFuncs([]utils.Middleware{middlewares.Log,}, controllers.SignUp)).Methods("POST").Name("SignUp")

	return nil
}
