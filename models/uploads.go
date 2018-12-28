package models

import (
	"fmt"
)

func (db *DB) SaveUpload(upload Upload) error {
	err := db.Set("glacier_uploads", upload.Metadata.Title, upload)
	//fmt.Printf("glacier upload saved with id: %d", upload.ArchiveCreationOutput.ArchiveId)
	return err
}

func (db *DB) AllUploads() (archiveCreationOutputs []Upload, err error) {
	fmt.Println("Fetching uploads")
	err = db.All(&archiveCreationOutputs)
	return archiveCreationOutputs, err
}
