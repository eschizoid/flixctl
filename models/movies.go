package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jrudio/go-plex-client" //nolint:goimports
)

type Movie struct {
	Metadata  plex.Metadata `json:"metadata"`
	Unwatched int           `json:"unwatched"`
	Title     string        `json:"title"`
}

func SavePlexMovie(movie Movie, svc *dynamodb.DynamoDB) (err error) {
	var av map[string]*dynamodb.AttributeValue
	av, _ = dynamodbattribute.MarshalMap(movie)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(plexMoviesTableName),
	}
	_, err = svc.PutItem(input)
	return err
}

func AllPlexMovies(svc *dynamodb.DynamoDB) (movies []Movie, err error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(plexMoviesTableName),
	}
	result, err := svc.Scan(params)
	for _, i := range result.Items {
		item := Movie{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
		}
		movies = append(movies, item)
	}
	return movies, err
}
