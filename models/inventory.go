package models

type InventoryArchive struct {
	ArchiveDescription string
	ArchiveID          string
	CreationDate       string
	SHA256TreeHash     string
	Size               int
}

func (db *DB) SaveInventoryArchive(inventoryArchive InventoryArchive) error {
	err := db.Set(inventoryArchivesBucketName, inventoryArchive.ArchiveID, inventoryArchive)
	return err
}

func (db *DB) AllInventoryArchives(keys [][]byte) (inventoryArchives []InventoryArchive, err error) {
	for _, key := range keys {
		if stringKey := string(key); stringKey != stormMetadataKey {
			var inventoryArchive InventoryArchive
			err = db.Get(inventoryArchivesBucketName, stringKey, &inventoryArchive)
			inventoryArchives = append(inventoryArchives, inventoryArchive)
		}
	}
	return inventoryArchives, err
}

func (db *DB) DeleteInventoryArchive(key string) (err error) {
	err = db.Delete(inventoryArchivesBucketName, key)
	return err
}

func (db *DB) DeleteAllInventoryArchives(keys [][]byte) (err error) {
	for _, key := range keys {
		if stringKey := string(key); stringKey != stormMetadataKey {
			err = db.Delete(inventoryArchivesBucketName, stringKey)
		}
	}
	return err
}
