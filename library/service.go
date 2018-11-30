package library

import (
	"os"

	"github.com/eschizoid/flixctl/models"
	"github.com/jrudio/go-plex-client"
	"go.etcd.io/bbolt"
)

var (
	DB        = models.NewDB(os.Getenv("BOLT_DATABASE"), []string{"plex_movies", "glacier_uploads"})
	PlexToken = os.Getenv("PLEX_TOKEN")
)

func GetLivePlexMovies(filter string) ([]plex.Metadata, error) {
	plexClient, err := plex.New("https://marianoflix.duckdns.org:32400", PlexToken)
	showError(err)
	var libraries plex.LibrarySections
	libraries, err = plexClient.GetLibraries()
	showError(err)
	directories := libraries.MediaContainer.Directory
	moviesDirectory := chooseMovies(directories, func(statusMessage string) bool { return statusMessage == "movie" })
	searchResults, err := plexClient.GetLibraryContent(moviesDirectory[0].Key, filter)
	showError(err)
	movies := searchResults.MediaContainer.Metadata
	return movies, err
}

func GetCachedPlexMovies() ([]plex.Metadata, error) {
	var err error
	var movies []plex.Metadata //nolint:prealloc
	var movie plex.Metadata
	keys := getAllKeys([]byte("plex_movies"))
	for _, key := range keys {
		err = DB.Get("plex_movies", string(key), &movie)
		showError(err)
		movies = append(movies, movie)
	}
	return movies, err
}

func SavePlexMovie(movie plex.Metadata) error {
	err := DB.SaveMovie(movie)
	return err
}

func GetGlacierMovies() ([]models.Upload, error) {
	var err error
	var uploads []models.Upload //nolint:prealloc
	var upload models.Upload
	keys := getAllKeys([]byte("glacier_uploads"))
	for _, key := range keys {
		err = DB.Get("glacier_uploads", string(key), &upload)
		showError(err)
		uploads = append(uploads, upload)
	}
	return uploads, err
}

func SaveGlacierMovie(upload models.Upload) error {
	err := DB.SaveUpload(upload)
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
