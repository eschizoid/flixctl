package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type OauthToken struct {
	ClientID string `json:"client_id"`
	Token    string `json:"token"`
}

func SaveOauthToken(clientID string, token string, svc *dynamodb.DynamoDB) (err error) {
	oauthToken := OauthToken{
		ClientID: clientID,
		Token:    token,
	}
	var av map[string]*dynamodb.AttributeValue
	av, _ = dynamodbattribute.MarshalMap(oauthToken)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(oauthTokenTableName),
	}
	_, err = svc.PutItem(input)
	return err
}
