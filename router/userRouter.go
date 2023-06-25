package router

import (
	"gRPC2/client/user"
	"gRPC2/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func UserRouter(r *mux.Router) {
	// Enable CORS middleware
	r.Use(middleware.EnableCors)
	// apply jwt Authentication
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Use(middleware.Auth())
	//client2 for register user microservice
	users := r.PathPrefix("/user").Subrouter()
	users.Path("/register").HandlerFunc(user.RegisterUser).Methods(http.MethodPost)
	users.Path("/login").HandlerFunc(user.LoginUsers).Methods(http.MethodPost)
	users.Path("/{user_id}").HandlerFunc(user.GetUser).Methods(http.MethodGet)
	users.Path("/").HandlerFunc(user.GetUsers).Methods(http.MethodGet)
	users.Path("/{user_id}").HandlerFunc(user.UpdateUser).Methods(http.MethodPut)
	users.Path("/{user_id}").HandlerFunc(user.DeleteUser).Methods(http.MethodDelete)

}
