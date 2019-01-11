package models

import (
	"github.com/jrudio/go-plex-client" //nolint:goimports
)

type Movie struct {
	Metadata  plex.Metadata
	Unwatched int
}

const plexMoviesBucketName = "plex_movies"

func (db *DB) SavePlexMovie(movie Movie) error {
	err := db.Set(plexMoviesBucketName, movie.Metadata.Title, movie)
	return err
}

func (db *DB) AllPlexMovies(keys [][]byte) (movies []Movie, err error) {
	for _, key := range keys {
		if stringKey := string(key); stringKey != StormMetadataKey {
			var movie Movie
			err = db.Get(plexMoviesBucketName, stringKey, &movie)
			movies = append(movies, movie)
		}
	}
	return movies, err
}
