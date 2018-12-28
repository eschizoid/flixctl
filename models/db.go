package models

import (
	"time"

	"github.com/asdine/storm"
	"github.com/jrudio/go-plex-client"
	"go.etcd.io/bbolt"
)

type Datastore interface {
	AllMovies() ([]plex.Metadata, error)
	SaveMovie(plex.Metadata) error
	AllUpload() ([]Upload, error)
	SaveUpload(Upload) error
	AllArchives() ([]Archive, error)
	SaveArchive(Archive) error
}

type DB struct {
	*storm.DB
}

func NewDB(dataSourceName string, buckets []string) *DB {
	db, _ := storm.Open(dataSourceName, storm.BoltOptions(0600, &bolt.Options{Timeout: 10 * time.Second}))
	db.Bolt.Update(func(tx *bolt.Tx) error { //nolint:errcheck
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
