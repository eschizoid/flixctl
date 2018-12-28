package models

import (
	"fmt"
)

func (db *DB) SaveInventoryArchive(inventoryArchive InventoryArchive) error {
	err := db.Set("glacier_archives", inventoryArchive.ArchiveID, inventoryArchive)
	//fmt.Printf("'%s' archive saved\n", archive.ArchiveID)
	return err
}

func (db *DB) AllInventoryArchives() (inventoryArchives []InventoryArchive, err error) {
	fmt.Println("Fetching glacier archives")
	err = db.All(&inventoryArchives)
	return inventoryArchives, err
}
