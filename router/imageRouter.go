package router

import (
	"gRPC2/client/image"
	"gRPC2/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func ImageRouter(r *mux.Router) {
	// Enable CORS middleware
	r.Use(middleware.EnableCors)
	// apply jwt Authentication
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Use(middleware.Auth())

	//client3 for upload image on aws microservice
	images := r.PathPrefix("/image").Subrouter()
	images.Path("/aws").HandlerFunc(image.UploadImage).Methods(http.MethodPost)
}
