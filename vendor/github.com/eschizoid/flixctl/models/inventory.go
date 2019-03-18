package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type InventoryArchive struct {
	ArchiveDescription string `json:"archive_description"`
	ArchiveID          string `json:"archive_id"`
	CreationDate       string `json:"creation_date"`
	SHA256TreeHash     string `json:"sha256_tree_hash"`
	Size               int    `json:"size"`
}

func SaveInventoryArchive(inventoryArchive InventoryArchive, svc *dynamodb.DynamoDB) (err error) {
	var av map[string]*dynamodb.AttributeValue
	av, _ = dynamodbattribute.MarshalMap(inventoryArchive)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(inventoryArchivesTableName),
	}
	_, err = svc.PutItem(input)
	return err
}

func AllInventoryArchives(svc *dynamodb.DynamoDB) (inventoryArchives []InventoryArchive, err error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(inventoryArchivesTableName),
	}
	result, _ := svc.Scan(params)
	for _, i := range result.Items {
		item := InventoryArchive{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}
		inventoryArchives = append(inventoryArchives, item)
	}
	return inventoryArchives, err
}

func DeleteInventoryArchive(key string, svc *dynamodb.DynamoDB) (err error) {
	result, _ := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(inventoryArchivesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"archive_id": {
				N: aws.String(key),
			},
		},
	})
	av, _ := dynamodbattribute.MarshalMap(result.Item)
	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: aws.String(inventoryArchivesTableName),
	}
	_, err = svc.DeleteItem(input)
	if err != nil {
		fmt.Println("Got error calling DeleteItem")
		fmt.Println(err.Error())
	}
	return err
}

func DeleteAllInventoryArchives(svc *dynamodb.DynamoDB) (err error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(inventoryArchivesTableName),
	}
	result, err := svc.Scan(params)
	for _, item := range result.Items {
		var av map[string]*dynamodb.AttributeValue
		av, _ = dynamodbattribute.MarshalMap(item)
		if err != nil {
			fmt.Println("Got error marshalling map:")
			fmt.Println(err.Error())
		}
		input := &dynamodb.DeleteItemInput{
			Key:       av,
			TableName: aws.String(inventoryArchivesTableName),
		}
		_, err = svc.DeleteItem(input)
		if err != nil {
			fmt.Println("Got error calling DeleteItem")
			fmt.Println(err.Error())
		}
	}
	return err
}
