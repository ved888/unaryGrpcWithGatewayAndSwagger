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

	//client3 for upload image on aws microservice
	images := r.PathPrefix("/image").Subrouter()
	// apply jwt Authentication
	images.Use(middleware.Auth())
	images.Path("/aws").HandlerFunc(image.UploadImage).Methods(http.MethodPost)
}
