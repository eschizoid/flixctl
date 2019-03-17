package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func executePlexCommand(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) error {
	//plex.Monitor()
	return nil
}

func main() {
	lambda.Start(executePlexCommand)
}
