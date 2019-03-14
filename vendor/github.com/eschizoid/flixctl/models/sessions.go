package models

import "time"

func (db *DB) SaveLastActiveSession(time time.Time) error {
	err := db.Set(plexSessionsBucketName, "last_activity", time)
	return err
}

func (db *DB) GetLastActiveSession() (lastActiveTime time.Time, err error) {
	err = db.Get(plexSessionsBucketName, "last_activity", &lastActiveTime)
	return lastActiveTime, err
}
