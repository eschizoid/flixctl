package models

func (db *DB) SaveOauthToken(clientID string, token string) error {
	err := db.Set(oauthTokenBucketName, clientID, token)
	return err
}
