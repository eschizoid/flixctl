package models

import (
	"fmt"
)

func (db *DB) SaveMovie(movie Movie) error {
	err := db.Set("plex_movies", movie.Metadata.Title, movie)
	fmt.Printf("'%s' movie saved\n", movie.Metadata.Title)
	return err
}

func (db *DB) AllMovies() (directories []Movie, err error) {
	fmt.Println("Fetching movies")
	err = db.All(&directories)
	return directories, err
}
