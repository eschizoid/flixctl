package models

import (
	"github.com/jrudio/go-plex-client"
)

type Movie struct {
	Metadata  plex.Metadata
	Unwatched int
}

const plexMoviesBucketName = "plex_movies"

func (db *DB) SavePlexMovie(movie Movie) error {
	err := db.Set(plexMoviesBucketName, movie.Metadata.Title, movie)
	//fmt.Printf("'%s' movie saved\n", movie.Metadata.Title)
	return err
}

func (db *DB) AllPlexMovies(keys [][]byte) (movies []Movie, err error) {
	//fmt.Println("Fetching movies")
	for _, key := range keys {
		var movie Movie
		err = db.Get(plexMoviesBucketName, key, &movie)
		movies = append(movies, movie)
	}
	return movies, err
}
