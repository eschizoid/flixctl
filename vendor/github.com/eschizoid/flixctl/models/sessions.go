package models

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type PlexSession struct {
	ID   string    `json:"id"`
	Time time.Time `json:"time"`
}

func SaveLastActiveSession(key string, svc *dynamodb.DynamoDB) (err error) {
	plexSession := PlexSession{
		ID:   key,
		Time: getTime(),
	}
	var av map[string]*dynamodb.AttributeValue
	av, _ = dynamodbattribute.MarshalMap(plexSession)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(plexSessionsTableName),
	}
	_, err = svc.PutItem(input)
	return err
}

func GetLastActiveSession(key string, svc *dynamodb.DynamoDB) (time.Time, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(plexSessionsTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(key),
					},
				},
			},
		},
	}
	var result, err = svc.Query(queryInput)
	item := PlexSession{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		fmt.Println("Got error unmarshalling:")
		fmt.Println(err.Error())
	}
	return item.Time, err
}

func getTime() time.Time {
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println(err)
	}
	return time.Now().In(location)
}
