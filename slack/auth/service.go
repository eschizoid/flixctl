package auth

import (
	"github.com/eschizoid/flixctl/models"
)

func SaveToken(clientID string, token string) error {
	err :=  models.Database.SaveOauthToken(clientID, token)
	return err
}
