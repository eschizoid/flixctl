package library

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eschizoid/flixctl/models"
	"github.com/jrudio/go-plex-client"
)

var (
	PlexToken = os.Getenv("PLEX_TOKEN")
)

func GetLivePlexMovies(filter int) (movies []models.Movie, err error) {
	var plexClient *plex.Plex
	plexClient, err = plex.New(fmt.Sprintf("https://%s:32400", os.Getenv("FLIXCTL_HOST")), PlexToken)
	showError(err)
	var libraries plex.LibrarySections
	libraries, err = plexClient.GetLibraries()
	showError(err)
	directories := libraries.MediaContainer.Directory
	moviesDirectory := findMoviesDirectory(directories, func(directoryType string) bool { return directoryType == "movie" })
	searchResults, err := plexClient.GetLibraryContent(moviesDirectory.Key, fmt.Sprintf("?unwatched=%d", filter))
	showError(err)
	for _, metadata := range searchResults.MediaContainer.Metadata {
		movie := models.Movie{
			Metadata:  metadata,
			Title:     metadata.Title,
			Unwatched: filter,
		}
		movies = append(movies, movie)
	}
	return movies, err
}

func DeleteteGlacierInventoryArchives(svc *dynamodb.DynamoDB) (err error) {
	return models.DeleteAllInventoryArchives(svc)
}

func DeleteteGlacierInventoryArchive(key string, svc *dynamodb.DynamoDB) (err error) {
	return models.DeleteInventoryArchive(key, svc)
}

func GetCachedPlexMovies(svc *dynamodb.DynamoDB) ([]models.Movie, error) {
	return models.AllPlexMovies(svc)
}

func GetGlacierMovies(svc *dynamodb.DynamoDB) (uploads []models.Upload, err error) {
	return models.AllUploads(svc)
}

func GetGlacierInventoryArchives(svc *dynamodb.DynamoDB) (archives []models.InventoryArchive, err error) {
	return models.AllInventoryArchives(svc)
}

func SaveGlacierInventoryArchive(archive models.InventoryArchive, svc *dynamodb.DynamoDB) error {
	err := models.SaveInventoryArchive(archive, svc)
	return err
}

func SaveGlacierMovie(upload models.Upload, svc *dynamodb.DynamoDB) error {
	err := models.SaveUpload(upload, svc)
	return err
}

func SavePlexMovie(movie models.Movie, svc *dynamodb.DynamoDB) error {
	err := models.SavePlexMovie(movie, svc)
	return err
}

func findMoviesDirectory(directories []plex.Directory, test func(string) bool) (movieDirectory plex.Directory) {
	for _, directory := range directories {
		if test(directory.Type) {
			movieDirectory = directory
			break
		}
	}
	return movieDirectory
}

func showError(err error) {
	if err != nil {
		panic(err)
	}
}
