package library

import (
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/eschizoid/flixctl/models"
	"github.com/jrudio/go-plex-client"
)

var db = models.NewDB("/tmp/library.db")

func GetPlexMovies(token string) (results plex.SearchResults) {
	plexClient, err := plex.New("https://marianoflix.duckdns.org:32400", token)
	if err != nil {
		panic(err)
	}
	libraries, err := plexClient.GetLibraries()
	if err != nil {
		panic(err)
	}
	directories := libraries.MediaContainer.Directory
	moviesDirectory := chooseMovies(directories, func(statusMessage string) bool { return statusMessage == "movie" })
	movies, err := plexClient.GetLibraryContent(moviesDirectory[0].Key, "")
	if err != nil {
		panic(err)
	}
	return movies
}

func SaveMovie(movie plex.Metadata) error {
	err := db.SaveMovie(movie)
	return err
}

func SaveUpload(archiveCreationOutput glacier.ArchiveCreationOutput) error {
	err := db.SaveUpload(archiveCreationOutput)
	return err
}

func GetGlacierMovies() (directories []plex.Directory, err error) {
	err = db.All(directories)
	return directories, err
}

func chooseMovies(directories []plex.Directory, test func(string) bool) (movieDirectory []plex.Directory) {
	for _, directory := range directories {
		if test(directory.Type) {
			movieDirectory = append(movieDirectory, directory)
		}
	}
	return
}
