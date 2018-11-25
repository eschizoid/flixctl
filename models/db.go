package models

import (
	"fmt"
	"os"
	"time"

	"github.com/asdine/storm"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/jrudio/go-plex-client"
	"go.etcd.io/bbolt"
)

type Datastore interface {
	AllMovies() ([]plex.Metadata, error)
	SaveMovie(plex.Metadata) error
	AllUpload() ([]glacier.ArchiveCreationOutput, error)
	SaveUpload(glacier.ArchiveCreationOutput) error
}

type DB struct {
	*storm.DB
}

func NewDB(dataSourceName string) *DB {
	pwd, _ := os.Getwd()
	db, _ := storm.Open(fmt.Sprintf("%s/database/storm/library.db", pwd), storm.BoltOptions(0600, &bolt.Options{Timeout: 10 * time.Second}))
	return &DB{db}
}
