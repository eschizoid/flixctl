package auth

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eschizoid/flixctl/models"
)

func SaveToken(clientID string, token string, svc *dynamodb.DynamoDB) error {
	err := models.SaveOauthToken(clientID, token, svc)
	return err
}
