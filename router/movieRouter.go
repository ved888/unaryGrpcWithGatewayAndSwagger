package router

import (
	"gRPC2/client/movie"
	"gRPC2/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func MovieRouter(r *mux.Router) {
	// Enable CORS middleware
	r.Use(middleware.EnableCors)

	//client for movie microservice
	movies := r.PathPrefix("/movies").Subrouter()
	movies.Use(middleware.Auth())
	movies.Path("/").HandlerFunc(movie.CreateMovie).Methods(http.MethodPost)
	movies.Path("/{id}").HandlerFunc(movie.GetMovieById).Methods(http.MethodGet)
	movies.Path("/").HandlerFunc(movie.GetAllMovie).Methods(http.MethodGet)
	movies.Path("/{id}").HandlerFunc(movie.UpdateMovieById).Methods(http.MethodPut)
	movies.Path("/{id}").HandlerFunc(movie.DeleteMovieById).Methods(http.MethodDelete)

}
