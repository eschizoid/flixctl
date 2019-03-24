package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func processDownloadsAndRequests(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) error {
	// get messages from both topics Downloads / Requests
	return nil
}

func main() {
	lambda.Start(processDownloadsAndRequests)
}
