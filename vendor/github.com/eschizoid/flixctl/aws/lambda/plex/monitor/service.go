package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/cmd/plex"
)

func executePlexCommand(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) error {
	fmt.Println(cloudWatchEvent.Time)
	plex.Monitor("false")
	return nil
}

func main() {
	lambda.Start(executePlexCommand)
}
