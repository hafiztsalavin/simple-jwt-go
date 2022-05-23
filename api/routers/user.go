package routers

import (
	"simple-jwt-go/api/controllers"
	"simple-jwt-go/api/middlewares"
	"simple-jwt-go/api/utils"

	"github.com/gorilla/mux"
)

// Define routes
func UserRoutes(router *mux.Router) error {
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST").Name("SignUp")
	router.HandleFunc("/signin", controllers.SignIn).Methods("POST").Name("SignIn")

	router.HandleFunc("/logout", utils.HandlerFuncs([]utils.Middleware{middlewares.CheckCookie}, controllers.Logout)).Methods("GET").Name("Logout")

	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.CheckCookie}, controllers.GetUser)).Methods("GET").Name("GetUser")
	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.CheckCookie}, controllers.DeleteUser)).Methods("DELETE").Name("DeleteUser")
	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.CheckCookie}, controllers.UpdateUser)).Methods("PATCH").Name("UpdateUser")

	return nil
}
