package main

import (
	"context"
	"flag"
	"fmt"
	"gRPC2/database"
	"gRPC2/dbhelper"
	"gRPC2/model"
	"gRPC2/pb/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type server struct {
	DAO dbhelper.DAO
	pb.UnimplementedMovieServiceServer
}

func (s *server) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	fmt.Println("Create Movie")
	movie := req.GetMovie()
	movie.Id = uuid.New().String()

	data := model.Movie{
		ID:    movie.GetId(),
		Title: movie.GetTitle(),
		Genre: movie.GetGenre(),
	}

	movieId, err := s.DAO.MovieCreate(&data)
	if err != nil {
		log.Printf("Failed to create movie: %v", err)
		return nil, err
	}

	return &pb.CreateMovieResponse{
		Movie: &pb.Movie{
			Id:    *movieId,
			Title: movie.GetTitle(),
			Genre: movie.GetGenre(),
		},
	}, nil
}

func (s *server) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	fmt.Println("Read Movie by id :", req.GetId())
	movie, err := s.DAO.GetMovieById(req.GetId())
	if err != nil {
		log.Printf("Failed to get movie: %v", err)
		return nil, err
	}
	return &pb.GetMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		},
	}, nil
}

func (s *server) GetMovies(ctx context.Context, req *pb.GetMoviesRequest) (*pb.GetMoviesResponse, error) {
	fmt.Println("Read Movies")
	movies, err := s.DAO.GetAllMovies()
	if err != nil {
		log.Printf("Failed to get movies: %v", err)
		return nil, err
	}
	pbMovies := make([]*pb.Movie, len(movies))
	for i, movie := range movies {
		pbMovies[i] = &pb.Movie{
			Id:    movie.ID,
			Title: movie.Title,
			Genre: movie.Genre,
		}
	}

	return &pb.GetMoviesResponse{
		Movies: pbMovies,
	}, nil
}

func (s *server) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	fmt.Println("update movie")
	movie := req.GetMovie()
	data := model.Movie{
		ID:    movie.GetId(),
		Title: movie.GetTitle(),
		Genre: movie.GetGenre(),
	}
	err := s.DAO.UpdateMovie(data, movie.GetId())
	if err != nil {
		log.Printf("Failed to update movie: %v", err)
		return nil, err
	}
	return &pb.UpdateMovieResponse{
		Movie: &pb.Movie{
			Id:    movie.GetId(),
			Title: movie.GetTitle(),
			Genre: movie.GetGenre(),
		},
	}, nil

}

func (s *server) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	fmt.Println("delete movie")
	err := s.DAO.DeleteMovie(req.GetId())
	if err != nil {
		log.Printf("Failed to delete movie: %v", err)
		return nil, err
	}
	return &pb.DeleteMovieResponse{
			Success: true,
		},
		nil
}
func main() {
	flag.Parse()

	db, err := database.DbConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()

	fmt.Println("gRPC server running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterMovieServiceServer(s, &server{DAO: dbhelper.DAO{
		DB: db,
	},
	})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
