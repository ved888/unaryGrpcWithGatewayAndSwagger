package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	user3 "gRPC2/client/image"
	user2 "gRPC2/client/user"
	"gRPC2/jwtUser"
	"gRPC2/model"
	"gRPC2/pb/pb"
	"github.com/ghodss/yaml"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {

	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMovieServiceClient(conn)

	// Read the Swagger JSON file
	swaggerJSON, err := ioutil.ReadFile("doc/swagger/apidocs.swagger.json")
	if err != nil {
		log.Fatalf("Failed to read Swagger JSON file: %v", err)
	}
	// Parse the JSON
	var swaggerData interface{}
	err = json.Unmarshal(swaggerJSON, &swaggerData)
	if err != nil {
		log.Fatalf("Failed to parse Swagger JSON: %v", err)
	}
	// Convert to YAML
	swaggerYAML, err := yaml.Marshal(swaggerData)
	if err != nil {
		log.Fatalf("Failed to convert Swagger JSON to YAML: %v", err)
	}
	// Write the YAML to a file
	err = ioutil.WriteFile("swagger.yaml", swaggerYAML, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to write Swagger YAML file: %v", err)
	}
	log.Println("Swagger JSON converted to YAML successfully!")

	r := mux.NewRouter()

	// set swagger ui router
	openAPIHandler := openapiv2.NewHandler()
	r.PathPrefix("/q/").Handler(openAPIHandler)

	// run swagger ui on browser by this url
	fmt.Println("Run swagger ui on browser by this url : http://localhost:8080/q/swagger-ui/")

	// Enable CORS middleware
	r.Use(enableCors)

	// apply jwt Authentication
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.Use(jwtUser.Auth())

	// Create movie by given fields
	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods(http.MethodPost)

	// Read movie by movie id
	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods(http.MethodGet)

	// Read all movies
	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods(http.MethodGet)

	// Update the moves by id
	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods(http.MethodPut)

	// delete movie by id
	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods(http.MethodDelete)

	//client2 for register user microservice
	user := r.PathPrefix("/user").Subrouter()
	user.Path("/register").HandlerFunc(user2.RegisterUser).Methods(http.MethodPost)
	user.Path("/login").HandlerFunc(user2.LoginUsers).Methods(http.MethodPost)
	user.Path("/{user_id}").HandlerFunc(user2.GetUser).Methods(http.MethodGet)
	user.Path("").HandlerFunc(user2.GetUsers).Methods(http.MethodGet)
	user.Path("/{user_id}/update").HandlerFunc(user2.UpdateUser).Methods(http.MethodPut)
	user.Path("/{user_id}/delete").HandlerFunc(user2.DeleteUser).Methods(http.MethodDelete)

	//client3 for upload image on aws microservice
	image := r.PathPrefix("/image").Subrouter()
	image.Path("/aws").HandlerFunc(user3.UploadImage).Methods(http.MethodPost)

	fmt.Println("http server client Running port:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}

// Enable CORS middleware
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Allow credentials (if needed)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			// Handle preflight requests
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
