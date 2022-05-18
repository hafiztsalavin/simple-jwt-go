package routers

import (
	"simple-jwt-go/api/controllers"
	"simple-jwt-go/api/middlewares"
	"simple-jwt-go/api/utils"

	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router) error {
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST").Name("SignUp")
	router.HandleFunc("/signin", controllers.SignIn).Methods("POST").Name("SignIn")

	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.CheckJWT}, controllers.GetUser)).Methods("GET").Name("GetUser")
	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.CheckJWT}, controllers.DeleteUser)).Methods("DELETE").Name("DeleteUser")
	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.CheckJWT}, controllers.UpdateUser)).Methods("PATCH").Name("UpdateUser")

	return nil
}
