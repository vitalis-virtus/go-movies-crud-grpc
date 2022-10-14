package models

// models package-

import (
	"github.com/jinzhu/gorm"
	"moviesapp.com/grpc/server/pkg/config"
)

var db *gorm.DB

type Movie struct {
	gorm.Model
	Id    string `json:"id"`
	Isbn  string `json:"isbn"`
	Title string `json:"title"`
}

func init() {
	config.Connect()
	db = config.GetDb()
	db.AutoMigrate(&Movie{})
}

// create new movie in db
func (m *Movie) CreateMovie() *Movie {
	db.NewRecord(m)
	db.Create(&m)
	return m
}

// get all movies in db
func GetMovies() []Movie {
	var Movies []Movie
	db.Find(&Movies)
	return Movies
}

// get movie by id in db
func GetMovieById(ID string) (*Movie, *gorm.DB) {
	var getMovie Movie
	db := db.Where("ID=?", ID).Find((&getMovie))
	return &getMovie, db
}

// delete movie by id from db
func DeleteMovie(ID string) Movie {
	var movie Movie
	db.Where("ID=?", ID).Delete(movie)
	return movie
}
