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
	//fmt.Printf("'%s' archive saved\n", archive.ArchiveID)
	return err
}

func (db *DB) AllInventoryArchives(keys [][]byte) (inventoryArchives []InventoryArchive, err error) {
	//	fmt.Println("Fetching glacier archives")
	for _, key := range keys {
		var inventoryArchive InventoryArchive
		err = db.Get(inventoryArchivesBucketName, string(key), &inventoryArchive)
		inventoryArchives = append(inventoryArchives, inventoryArchive)
	}
	return inventoryArchives, err
}
