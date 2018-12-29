package models

func (db *DB) SavePlexMovie(movie Movie) error {
	err := db.Set("plex_movies", movie.Metadata.Title, movie)
	//fmt.Printf("'%s' movie saved\n", movie.Metadata.Title)
	return err
}

func (db *DB) AllPlexMovies() (directories []Movie, err error) {
	//fmt.Println("Fetching movies")
	err = db.All(&directories)
	return directories, err
}
