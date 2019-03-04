package auth

import (
	"os"

	"github.com/eschizoid/flixctl/models"
)

const (
	oauthTokens = "oauth_tokens" //nolint:gosec
)

var (
	Database = models.NewDB(os.Getenv("BOLT_DATABASE"), []string{oauthTokens})
)

func SaveToken(clientID string, token string) error {
	err := Database.SaveOauthToken(clientID, token)
	return err
}
