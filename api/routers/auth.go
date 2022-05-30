package routers

import (
	"simple-jwt-go/api/controllers"
	"simple-jwt-go/api/middlewares"
	"simple-jwt-go/api/utils"

	"github.com/gorilla/mux"
)

// Define routes
func UserAuthRoutes(router *mux.Router) error {
	userRouter := router.PathPrefix("/auth").Subrouter()
	userRouter.Path("/profile").HandlerFunc(utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckAccess}, controllers.GetUser)).Methods("GET").Name("GetUser")
	userRouter.Path("/user").HandlerFunc(utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckAccess}, controllers.DeleteUser)).Methods("DELETE").Name("DeleteUser")
	userRouter.Path("/user").HandlerFunc(utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckAccess}, controllers.UpdateUser)).Methods("PATCH").Name("UpdateUser")
	userRouter.Path("/logout").HandlerFunc(utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckAccess}, controllers.Logout)).Methods("GET").Name("Logout")
	userRouter.Path("/refresh").HandlerFunc(utils.HandlerFuncs([]utils.Middleware{middlewares.Log, middlewares.CheckRefresh}, controllers.Refresh)).Methods("GET").Name("Refresh")

	return nil
}
