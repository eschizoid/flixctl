package library

import (
	"fmt"
	"os"

	"github.com/eschizoid/flixctl/models"
	"github.com/jrudio/go-plex-client" //nolint:goimports
	"go.etcd.io/bbolt"
)

const (
	inventoryArchivesBucketName = "glacier_archives"
	plexMoviesBucketName        = "plex_movies"
	uploadsBucketName           = "glacier_uploads"
)

var (
	Database  = models.NewDB(os.Getenv("BOLT_DATABASE"), []string{plexMoviesBucketName, uploadsBucketName, inventoryArchivesBucketName})
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

func DeleteteGlacierInventoryArchives() (err error) {
	keys := getAllKeys([]byte(inventoryArchivesBucketName))
	return Database.DeleteAllInventoryArchives(keys)
}

func DeleteteGlacierInventoryArchive(key string) (err error) {
	return Database.DeleteInventoryArchive(key)
}

func GetCachedPlexMovies() ([]models.Movie, error) {
	keys := getAllKeys([]byte(plexMoviesBucketName))
	return Database.AllPlexMovies(keys)
}

func GetGlacierMovies() (uploads []models.Upload, err error) {
	keys := getAllKeys([]byte(uploadsBucketName))
	return Database.AllUploads(keys)
}

func GetGlacierInventoryArchives() (archives []models.InventoryArchive, err error) {
	keys := getAllKeys([]byte(inventoryArchivesBucketName))
	return Database.AllInventoryArchives(keys)
}

func FindGlacierMovie(title string) (archives models.Upload, err error) {
	return Database.FindUploadByID(title)
}

func SaveGlacierInventoryArchive(archive models.InventoryArchive) error {
	err := Database.SaveInventoryArchive(archive)
	return err
}

func SaveGlacierMovie(upload models.Upload) error {
	err := Database.SaveUpload(upload)
	return err
}

func SavePlexMovie(movie models.Movie) error {
	err := Database.SavePlexMovie(movie)
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

func getAllKeys(bucket []byte) [][]byte {
	var keys [][]byte
	Database.Bolt.View(func(tx *bolt.Tx) error { //nolint:errcheck
		b := tx.Bucket(bucket)
		_ = b.ForEach(func(k, v []byte) error { //nolint:errcheck
			// Due to
			// Byte slices returned from Bolt are only valid during a transaction. Once the transaction has been committed or rolled back then the memory they point to can be reused by a new page or can be unmapped from virtual memory and you'll see an unexpected fault address panic when accessing it.
			// We copy the slice to retain it
			dst := make([]byte, len(k))
			copy(dst, k)
			keys = append(keys, dst)
			return nil
		})
		return nil
	})
	return keys
}

func showError(err error) {
	if err != nil {
		panic(err)
	}
}
