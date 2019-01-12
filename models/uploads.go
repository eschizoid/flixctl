package models

import (
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/jrudio/go-plex-client" //nolint:goimports
)

type Upload struct {
	Metadata              plex.Metadata
	ArchiveCreationOutput glacier.ArchiveCreationOutput
}

const uploadsBucketName = "glacier_uploads"

func (db *DB) SaveUpload(upload Upload) error {
	err := db.Set(uploadsBucketName, upload.Metadata.Title, upload)
	return err
}

func (db *DB) AllUploads(keys [][]byte) (uploads []Upload, err error) {
	for _, key := range keys {
		if stringKey := string(key); stringKey != StormMetadataKey {
			var upload Upload
			err = db.Get(uploadsBucketName, stringKey, &upload)
			uploads = append(uploads, upload)
		}
	}
	return uploads, err
}

func (db *DB) FindUploadByID(title string) (Upload, error) {
	var upload Upload
	err := db.Get(uploadsBucketName, title, &upload)
	return upload, err
}
