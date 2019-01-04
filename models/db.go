package models

import (
	"time"

	"github.com/asdine/storm"
	"github.com/jrudio/go-plex-client" //nolint:goimports
	"go.etcd.io/bbolt"
)

type Datastore interface {
	AllInventoryArchives([][]byte) ([]InventoryArchive, error)
	AllPlexMovies([][]byte) ([]plex.Metadata, error)
	AllUploads([][]byte) ([]Upload, error)
	DeleteAllInventoryArchives([][]byte) error
	FindUploadByID(string) (Upload, error)
	SaveInventoryArchive(InventoryArchive) error
	SavePlexMovie(plex.Metadata) error
	SaveUpload(Upload) error
}

type DB struct {
	*storm.DB
}

func NewDB(dataSourceName string, buckets []string) *DB {
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
