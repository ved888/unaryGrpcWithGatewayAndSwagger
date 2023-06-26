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

	//client2 for register user microservice
	users := r.PathPrefix("/user").Subrouter()
	users.Path("/register").HandlerFunc(user.RegisterUser).Methods(http.MethodPost)
	users.Path("/login").HandlerFunc(user.LoginUsers).Methods(http.MethodPost)
	userAction := users.PathPrefix("/auth").Subrouter()
	// apply jwt Authentication
	userAction.Use(middleware.Auth())
	userAction.Path("/{user_id}").HandlerFunc(user.GetUser).Methods(http.MethodGet)
	userAction.Path("/").HandlerFunc(user.GetUsers).Methods(http.MethodGet)
	userAction.Path("/{user_id}").HandlerFunc(user.UpdateUser).Methods(http.MethodPut)
	userAction.Path("/{user_id}").HandlerFunc(user.DeleteUser).Methods(http.MethodDelete)

}
