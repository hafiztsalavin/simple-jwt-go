package routers

import (
	"simple-jwt-go/api/controllers"
	"simple-jwt-go/api/middlewares"
	"simple-jwt-go/api/utils"

	"github.com/gorilla/mux"
)

// Define routes
func UserRoutes(router *mux.Router) error {
	router.HandleFunc("/register", utils.HandlerFuncs([]utils.Middleware{middlewares.Log}, controllers.SignUp)).Methods("POST").Name("SignUp")
	router.HandleFunc("/login", utils.HandlerFuncs([]utils.Middleware{middlewares.Log}, controllers.SignIn)).Methods("POST").Name("SignIn")
	router.HandleFunc("/logout", utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckAccess}, controllers.Logout)).Methods("GET").Name("Logout")
	router.HandleFunc("/refresh", utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckRefresh}, controllers.Refresh)).Methods("GET").Name("Refresh")

	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckAccess}, controllers.GetUser)).Methods("GET").Name("GetUser")
	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckAccess}, controllers.DeleteUser)).Methods("DELETE").Name("DeleteUser")
	router.HandleFunc("/account", utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckAccess}, controllers.UpdateUser)).Methods("PATCH").Name("UpdateUser")

	return nil
}
