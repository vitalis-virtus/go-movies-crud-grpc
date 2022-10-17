package main

// package main
// gRPC server

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	pb "moviesapp.com/grpc/protos"
	"moviesapp.com/grpc/server/pkg/models"
)

const (
	port = ":50051"
)

type movieServer struct {
	pb.UnimplementedMovieServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterMovieServer(s, &movieServer{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *movieServer) GetMovies(in *pb.Empty,
	stream pb.Movie_GetMoviesServer) error {
	log.Printf("Received GetAllMovies")
	movies := models.GetMovies()
	for _, movie := range movies {
		if err := stream.Send(&pb.MovieInfo{Id: movie.GetMovieId(), Isbn: movie.Isbn, Title: movie.Title}); err != nil {
			return err
		}
	}
	return nil

}

func (s *movieServer) GetMovie(ctx context.Context,
	in *pb.Id) (*pb.MovieInfo, error) {
	log.Printf("Received GetMovie with id: %v", in)

	ID := in.GetValue()

	movieDetails, _ := models.GetMovieById(ID)

	if movieDetails.GetMovieId() == "" {
		return nil, errors.New("No movie with such id")
	}

	res := &pb.MovieInfo{Id: movieDetails.GetMovieId(), Isbn: movieDetails.Isbn, Title: movieDetails.Title}

	return res, nil
}

func (s *movieServer) CreateMovie(ctx context.Context,
	in *pb.MovieInfo) (*pb.Id, error) {
	log.Printf("Received CreateMovie with MovieInfo: %v", in)

	ID := pb.Id{}

	rand.Seed(time.Now().UnixNano())

	for {
		ID.Value = strconv.Itoa(rand.Intn(10000000000))

		movieDetails, _ := models.GetMovieById(ID.Value)

		if movieDetails.Id != ID.Value {
			break
		}

	}

	newMovie := &models.Movie{Id: ID.Value, Isbn: in.GetIsbn(), Title: in.GetTitle()}

	newMovie.CreateMovie()

	return &ID, nil
}

func (s *movieServer) UpdateMovie(ctx context.Context,
	in *pb.MovieInfo) (*pb.Status, error) {
	log.Printf("Received UpdateMovie with MovieInfo: %v", in)

	ID := in.GetId()

	movieDetails, db := models.GetMovieById(ID)

	// we are checking for updated fields from user and replace this fields in moviesDetails
	if in.GetIsbn() != "" {
		movieDetails.Isbn = in.GetIsbn()
	}
	if in.GetTitle() != "" {
		movieDetails.Title = in.GetTitle()
	}

	movieDetails.Id = ID

	db.Save(movieDetails)

	res := pb.Status{}

	res.Value = 1

	return &res, nil
}

func (s *movieServer) DeleteMovie(ctx context.Context,
	in *pb.Id) (*pb.Status, error) {
	log.Printf("Received DeleteMovie with id: %v", in)

	res := pb.Status{}

	ID := in.GetValue()

	models.DeleteMovie(ID)

	res.Value = 1

	return &res, nil
}
