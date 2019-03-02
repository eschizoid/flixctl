package models

const oauthTokenBucketName = "oauth_token" //nolint:gosec

func (db *DB) SaveOauthToken(clientID string, token string) error {
	err := db.Set(oauthTokenBucketName, clientID, token)
	return err
}
