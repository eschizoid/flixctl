package models

import (
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/jrudio/go-plex-client"
)

type Upload struct {
	Metadata              plex.Metadata
	ArchiveCreationOutput glacier.ArchiveCreationOutput
}

const uploadsBucketName = "glacier_uploads"

func (db *DB) SaveUpload(upload Upload) error {
	err := db.Set(uploadsBucketName, upload.Metadata.Title, upload)
	//fmt.Printf("glacier upload saved with id: %d", upload.ArchiveCreationOutput.ArchiveId)
	return err
}

func (db *DB) AllUploads(keys [][]byte) (uploads []Upload, err error) {
	//fmt.Println("Fetching uploads")
	for _, key := range keys {
		var upload Upload
		err = db.Get(uploadsBucketName, string(key), &upload)
		uploads = append(uploads, upload)
	}
	return uploads, err
}

func (db *DB) FindUploadByID(title string) (Upload, error) {
	var upload Upload
	err := db.Get(uploadsBucketName, title, &upload)
	return upload, err
}
