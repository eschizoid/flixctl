package library

import (
	"github.com/asdine/storm"
	"github.com/jrudio/go-plex-client"
)

func FindInPlex(query string) (results plex.SearchResults) {
	plexConnection, err := plex.New("http://192.168.1.2:32400", "myPlexToken")
	if err != nil {
		panic(err)
	}
	results, _ = plexConnection.Search(query)
	return results
}

func FindInLibrary(query string) (results plex.SearchResults) {
	db, err := storm.Open("library.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//db.Find()
	return *new(plex.SearchResults)
}
