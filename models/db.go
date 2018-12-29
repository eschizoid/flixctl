package models

import (
	"time"

	"github.com/asdine/storm"
	"github.com/jrudio/go-plex-client"
	"go.etcd.io/bbolt"
)

type Datastore interface {
	AllPlexMovies() ([]plex.Metadata, error)
	SavePlexMovie(plex.Metadata) error
	AllUploads() ([]Upload, error)
	SaveUpload(Upload) error
	AllInventoryArchives() ([]InventoryArchive, error)
	SaveInventoryArchive(InventoryArchive) error
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
