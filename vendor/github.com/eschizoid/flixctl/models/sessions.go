package models

import "time"

func (db *DB) SaveLastActiveSession(time time.Time) error {
	err := db.Set(plexSessionsBucketName, "last_activity", time.String())
	return err
}

func (db *DB) GetLastActiveSession() (time.Time, error) {
	var lastActiveTime time.Time
	err := db.Get(plexSessionsBucketName, "last_activity", &lastActiveTime)
	return lastActiveTime, err
}
