package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	pb "moviesapp.com/grpc/protos"
	"moviesapp.com/grpc/server/pkg/models"
	"net"
	"strconv"
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
	log.Printf("Received: %v", in)
	movies := models.GetMovies()
	for _, movie := range movies {
		if err := stream.Send(&pb.MovieInfo{Id: movie.Id, Isbn: movie.Isbn, Title: movie.Title}); err != nil {
			return err
		}
	}
	return nil

}

func (s *movieServer) GetMovie(ctx context.Context,
	in *pb.Id) (*pb.MovieInfo, error) {
	log.Printf("Received: %v", in)

	ID := in.GetValue()

	movieDetails, _ := models.GetMovieById(ID)

	res := &pb.MovieInfo{Id: movieDetails.Id, Isbn: movieDetails.Isbn, Title: movieDetails.Title}

	return res, nil
}

func (s *movieServer) CreateMovie(ctx context.Context,
	in *pb.MovieInfo) (*pb.Id, error) {
	log.Printf("Received: %v", in)

	ID := pb.Id{}
	ID.Value = strconv.Itoa(rand.Intn(100000000))

	newMovie := &models.Movie{Id: ID.Value, Isbn: in.GetIsbn(), Title: in.GetTitle()}

	newMovie.CreateMovie()

	return &ID, nil
}

func (s *movieServer) UpdateMovie(ctx context.Context,
	in *pb.MovieInfo) (*pb.Status, error) {
	log.Printf("Received: %v", in)

	ID := in.GetId()

	movieDetails, db := models.GetMovieById(ID)

	// we are checking for updated fields from user and replace this fields in moviesDetails
	if in.GetIsbn() != "" {
		movieDetails.Isbn = in.GetIsbn()
	}
	if in.GetTitle() != "" {
		movieDetails.Title = in.GetTitle()
	}

	db.Save(&movieDetails)

	res := pb.Status{}

	res.Value = 1

	return &res, nil
}

func (s *movieServer) DeleteMovie(ctx context.Context,
	in *pb.Id) (*pb.Status, error) {
	log.Printf("Received: %v", in)

	res := pb.Status{}

	ID := in.GetValue()

	models.DeleteMovie(ID)

	res.Value = 1

	return &res, nil
}
