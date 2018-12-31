package models

import (
	"github.com/jrudio/go-plex-client"
)

type Movie struct {
	Metadata  plex.Metadata
	Unwatched int
}

func (db *DB) SavePlexMovie(movie Movie) error {
	err := db.Set("plex_movies", movie.Metadata.Title, movie)
	//fmt.Printf("'%s' movie saved\n", movie.Metadata.Title)
	return err
}

func (db *DB) AllPlexMovies() (movies []Movie, err error) {
	//fmt.Println("Fetching movies")
	err = db.All(&movies)
	return movies, err
}
