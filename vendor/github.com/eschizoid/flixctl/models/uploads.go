package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/jrudio/go-plex-client"
)

type Upload struct {
	Metadata              *plex.Metadata
	ArchiveCreationOutput *glacier.ArchiveCreationOutput
}

func (db *DB) SaveUpload(upload Upload) error {
	err := db.Set("glacier_uploads", upload.ArchiveCreationOutput, upload)
	fmt.Printf("glacier upload saved wiht id: %d", upload.ArchiveCreationOutput.ArchiveId)
	return err
}

func (db *DB) AllUploads() (archiveCreationOutputs []Upload, err error) {
	fmt.Println("Fetching uploads")
	err = db.All(&archiveCreationOutputs)
	return archiveCreationOutputs, err
}
