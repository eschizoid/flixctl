package models

import (
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/jrudio/go-plex-client"
)

type Upload struct {
	Metadata              plex.Metadata
	ArchiveCreationOutput glacier.ArchiveCreationOutput
}

func (db *DB) SaveUpload(upload Upload) error {
	err := db.Set("glacier_uploads", upload.Metadata.Title, upload)
	//fmt.Printf("glacier upload saved with id: %d", upload.ArchiveCreationOutput.ArchiveId)
	return err
}

func (db *DB) AllUploads() (archiveCreationOutputs []Upload, err error) {
	//fmt.Println("Fetching uploads")
	err = db.All(&archiveCreationOutputs)
	return archiveCreationOutputs, err
}

func (db *DB) FindUploadByID(title string) (Upload, error) {
	var upload Upload
	err := db.Get("glacier_uploads", title, &upload)
	return upload, err
}
