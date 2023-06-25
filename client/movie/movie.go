package movie

import (
	"context"
	"encoding/json"
	"gRPC2/model"
	"gRPC2/pb/pb"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

var Conn, _ = grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

var client = pb.NewMovieServiceClient(Conn)

// CreateMovie
// @Summary create movie for given particular field.
// @Description create movie with the input payload.
// @Tags Movie
// @Param movie body model.Movie true "Create movie"
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Failure 500
// @Security     ApiKeyAuth
// @Router /movies/ [post]
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if movie.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if movie.Genre == "" {
		http.Error(w, "Genre is required", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	res, err := client.CreateMovie(ctx, &pb.CreateMovieRequest{
		Movie: &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res.Movie)
	if err != nil {
		return
	}
}

// GetMovieById
// @Summary Get details of movie by movie id
// @Description Get details of movie by id
// @Tags Movie
// @Param id path string true "get movie by id"
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400
// @Failure 500
// @Security ApiKeyAuth
// @Router /movies/{id} [get]
func GetMovieById(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie

	movie.ID = mux.Vars(r)["id"]
	if movie.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	res, err := client.GetMovie(ctx, &pb.GetMovieRequest{
		Id: movie.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res.Movie)
	if err != nil {
		return
	}
}

// GetAllMovie
// @Summary Get details of All movie
// @Description Get details of movies
// @Tags Movie
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400
// @Failure 500
// @Security ApiKeyAuth
// @Router /movies/ [get]
func GetAllMovie(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := client.GetMovies(ctx, &pb.GetMoviesRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res.Movies)
	if err != nil {
		return
	}
}

// UpdateMovieById
// @Summary update movie for given field.
// @Description update movie with the given movie fields.
// @Tags Movie
// @Param movie body model.Movie true "update movie"
// @Param id path string true "update movie by id"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 500
// @Security ApiKeyAuth
// @Router /movies/{id} [put]
func UpdateMovieById(w http.ResponseWriter, r *http.Request) {
	movieID := mux.Vars(r)["id"]

	var movie model.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if movie.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if movie.Genre == "" {
		http.Error(w, "genre is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	res, err := client.UpdateMovie(ctx, &pb.UpdateMovieRequest{
		Id: movieID,
		Movie: &pb.Movie{
			Id:    movieID,
			Title: movie.Title,
			Genre: movie.Genre,
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res.Movie)
	if err != nil {
		return
	}
}

// DeleteMovieById
// @Summary delete details of movie entry by movie id
// @Description delete details of movie by id
// @Tags Movie
// @Param id path string true "delete movie by id"
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400
// @Failure 500
// @Security ApiKeyAuth
// @Router /movies/{id} [delete]
func DeleteMovieById(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie

	movie.ID = mux.Vars(r)["id"]
	if movie.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	res, err := client.DeleteMovie(ctx, &pb.DeleteMovieRequest{
		Id: movie.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res.Success)
	if err != nil {
		return
	}
}
