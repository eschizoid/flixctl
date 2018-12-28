package library

import (
	"fmt"
	"os"

	"github.com/eschizoid/flixctl/models"
	"github.com/jrudio/go-plex-client"
	"go.etcd.io/bbolt"
)

var (
	DB        = models.NewDB(os.Getenv("BOLT_DATABASE"), []string{"plex_movies", "glacier_uploads", "glacier_archives"})
	PlexToken = os.Getenv("PLEX_TOKEN")
)

func GetLivePlexMovies(filter int) (movies []models.Movie, err error) {
	var plexClient *plex.Plex
	plexClient, err = plex.New("https://marianoflix.duckdns.org:32400", PlexToken)
	showError(err)
	var libraries plex.LibrarySections
	libraries, err = plexClient.GetLibraries()
	showError(err)
	directories := libraries.MediaContainer.Directory
	moviesDirectory := chooseMovies(directories, func(statusMessage string) bool { return statusMessage == "movie" })
	searchResults, err := plexClient.GetLibraryContent(moviesDirectory[0].Key, fmt.Sprintf("?unwatched=%d", filter))
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

func GetCachedPlexMovies() (movies []models.Movie, err error) {
	keys := getAllKeys([]byte("plex_movies"))
	for _, key := range keys {
		var movie models.Movie
		err = DB.Get("plex_movies", string(key), &movie)
		movies = append(movies, movie)
	}
	return movies, err
}

func GetGlacierMovies() (uploads []models.Upload, err error) {
	keys := getAllKeys([]byte("glacier_uploads"))
	for _, key := range keys {
		var upload models.Upload
		err = DB.Get("glacier_uploads", string(key), &upload)
		uploads = append(uploads, upload)
	}
	return uploads, err
}

func GetGlacierInventoryArchives() (archives []models.InventoryArchive, err error) {
	keys := getAllKeys([]byte("glacier_archives"))
	for _, key := range keys {
		var archive models.InventoryArchive
		err = DB.Get("glacier_archives", string(key), &archive)
		archives = append(archives, archive)
	}
	return archives, err
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

func getAllKeys(bucket []byte) [][]byte {
	var keys [][]byte
	DB.Bolt.View(func(tx *bolt.Tx) error { //nolint:errcheck
		b := tx.Bucket(bucket)
		b.ForEach(func(k, v []byte) error { //nolint:errcheck
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
	if len(keys) > 0 {
		keys = keys[:len(keys)-1]
	}
	return keys
}

func chooseMovies(directories []plex.Directory, test func(string) bool) (movieDirectory []plex.Directory) {
	for _, directory := range directories {
		if test(directory.Type) {
			movieDirectory = append(movieDirectory, directory)
		}
	}
	return
}

func showError(err error) {
	if err != nil {
		panic(err)
	}
}
