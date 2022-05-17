package routers

import (
	"simple-jwt-go/api/controllers"
	"simple-jwt-go/api/utils"

	"simple-jwt-go/api/middlewares"

	"github.com/gorilla/mux"
)

func AddUserRoutes(router *mux.Router) error {
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST").Name("SignUp")
	router.HandleFunc("/signin", controllers.SignIn).Methods("POST").Name("SignIn")

	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.Path("/{id}").HandlerFunc(utils.ChainHandlerFuncs([]utils.Middleware{middlewares.CheckJWT}, controllers.GetUser)).Methods("GET").Name("GetUser")
	userRouter.Queries().HandlerFunc(utils.ChainHandlerFuncs([]utils.Middleware{middlewares.CheckJWT}, controllers.UpdateUser)).Methods("PATCH").Name("UpdateUser")
	userRouter.Path("/{id}").HandlerFunc(utils.ChainHandlerFuncs([]utils.Middleware{middlewares.CheckJWT}, controllers.DeleteUser)).Methods("DELETE").Name("DeleteUser")

	return nil
}
