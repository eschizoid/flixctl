package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/jrudio/go-plex-client"
)

type Upload struct {
	SearchResults         *plex.SearchResults
	ArchiveCreationOutput *glacier.ArchiveCreationOutput
}

func (db *DB) SaveUpload(upload glacier.ArchiveCreationOutput) error {
	fmt.Printf("Saving upload wiht id: %d", upload.ArchiveId)
	err := db.Set("glacier_uploads", upload.ArchiveId, upload)
	return err
}

func (db *DB) AllUploads() (archiveCreationOutputs []glacier.ArchiveCreationOutput, err error) {
	fmt.Println("Fetching uploads")
	err = db.All(&archiveCreationOutputs)
	return archiveCreationOutputs, err
}
