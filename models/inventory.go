package models

type InventoryArchive struct {
	ArchiveDescription string
	ArchiveID          string
	CreationDate       string
	SHA256TreeHash     string
	Size               int
}

const inventoryArchivesBucketName = "glacier_archives"

func (db *DB) SaveInventoryArchive(inventoryArchive InventoryArchive) error {
	err := db.Set(inventoryArchivesBucketName, inventoryArchive.ArchiveID, inventoryArchive)
	return err
}

func (db *DB) AllInventoryArchives(keys [][]byte) (inventoryArchives []InventoryArchive, err error) {
	for _, key := range keys {
		var inventoryArchive InventoryArchive
		err = db.Get(inventoryArchivesBucketName, string(key), &inventoryArchive)
		inventoryArchives = append(inventoryArchives, inventoryArchive)
	}
	return inventoryArchives, err
}

func (db *DB) DeleteAllInventoryArchives(keys [][]byte) (err error) {
	for _, key := range keys {
		err = db.Delete(inventoryArchivesBucketName, string(key))
	}
	return err
}
