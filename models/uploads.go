package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/jrudio/go-plex-client"
)

type Upload struct {
	Metadata              plex.Metadata                 `json:"metadata"`
	ArchiveCreationOutput glacier.ArchiveCreationOutput `json:"archive_creation_output"`
	Title                 string                        `json:"title"`
}

func SaveUpload(upload Upload, svc *dynamodb.DynamoDB) (err error) {
	var av map[string]*dynamodb.AttributeValue
	av, _ = dynamodbattribute.MarshalMap(upload)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(uploadsTableName),
	}
	_, err = svc.PutItem(input)
	return err
}

func AllUploads(svc *dynamodb.DynamoDB) (uploads []Upload, err error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(uploadsTableName),
	}
	result, err := svc.Scan(params)
	for _, i := range result.Items {
		item := Upload{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}
		uploads = append(uploads, item)
	}
	return uploads, err
}

func FindUploadByID(title string, svc *dynamodb.DynamoDB) (Upload, error) {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(plexSessionsTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"title": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(title),
					},
				},
			},
		},
	}
	var result, err = svc.Query(queryInput)
	item := Upload{}
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		fmt.Println("Got error unmarshalling:")
		fmt.Println(err.Error())
	}
	return item, err
}
