package models

import (
	"os"
	"time"

	"github.com/asdine/storm"
	"github.com/jrudio/go-plex-client" //nolint:goimports
	bolt "go.etcd.io/bbolt"
)

const (
	inventoryArchivesBucketName = "glacier_archives"
	plexMoviesBucketName        = "plex_movies"
	plexSessionsBucketName      = "plex_sessions"
	uploadsBucketName           = "glacier_uploads"
	oauthTokenBucketName        = "oauth_tokens" //nolint:gosec
	stormMetadataKey            = "__storm_metadata"
)

type Datastore interface {
	AllInventoryArchives([][]byte) ([]InventoryArchive, error)
	AllPlexMovies([][]byte) ([]plex.Metadata, error)
	AllUploads([][]byte) ([]Upload, error)
	DeleteAllInventoryArchives([][]byte) error
	FindUploadByID(string) (Upload, error)
	SaveInventoryArchive(InventoryArchive) error
	SaveOauthToken(string, string) error
	SavePlexMovie(plex.Metadata) error
	SaveLastActiveSession(time.Time) error
	GetLastActiveSession(Upload) error
}

var (
	Database = NewDB(os.Getenv("BOLT_DATABASE"))
)

type DB struct {
	*storm.DB
}

func NewDB(dataSourceName string) *DB {
	buckets := []string{plexMoviesBucketName, uploadsBucketName, inventoryArchivesBucketName, oauthTokenBucketName, plexSessionsBucketName}
	db, _ := storm.Open(dataSourceName, storm.BoltOptions(0600, &bolt.Options{Timeout: 10 * time.Second}))
	_ = db.Bolt.Update(func(tx *bolt.Tx) error { //nolint:errcheck
		for _, value := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(value))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &DB{db}
}
