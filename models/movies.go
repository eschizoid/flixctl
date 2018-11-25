package models

import (
	"fmt"

	"github.com/jrudio/go-plex-client"
)

func (db *DB) SaveMovie(movie plex.Metadata) error {
	err := db.Set("plex_movies", movie.Title, movie)
	fmt.Printf("'%s' movie metadata saved\n", movie.Title)
	return err
}

func (db *DB) AllMovies() (directories []plex.Metadata, err error) {
	fmt.Println("Fetching movies")
	err = db.All(&directories)
	return directories, err
}
