package models

import (
	"fmt"
)

func (db *DB) SaveArchive(archive Archive) error {
	err := db.Set("glacier_archives", archive.ArchiveID, archive)
	//fmt.Printf("'%s' archive saved\n", archive.ArchiveID)
	return err
}

func (db *DB) AllArchives() (archives []Archive, err error) {
	fmt.Println("Fetching movies")
	err = db.All(&archives)
	return archives, err
}
