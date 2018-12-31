package library

import (
	"fmt"
	"os"

	"github.com/eschizoid/flixctl/models"
	"github.com/jrudio/go-plex-client"
)

var (
	DB        = models.NewDB(os.Getenv("BOLT_DATABASE"), []string{"plex_movies", "glacier_uploads", "glacier_archives"})
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
			Unwatched: filter,
		}
		movies = append(movies, movie)
	}
	return movies, err
}

func GetCachedPlexMovies() ([]models.Movie, error) {
	return DB.AllPlexMovies()
}

func GetGlacierMovies() (uploads []models.Upload, err error) {
	return DB.AllUploads()
}

func GetGlacierInventoryArchives() (archives []models.InventoryArchive, err error) {
	return DB.AllInventoryArchives()
}

func FindGlacierMovie(title string) (archives models.Upload, err error) {
	return DB.FindUploadByID(title)
}

func SaveGlacierInventoryArchive(archive models.InventoryArchive) error {
	err := DB.SaveInventoryArchive(archive)
	return err
}

func SaveGlacierMovie(upload models.Upload) error {
	err := DB.SaveUpload(upload)
	return err
}

func SavePlexMovie(movie models.Movie) error {
	err := DB.SavePlexMovie(movie)
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
