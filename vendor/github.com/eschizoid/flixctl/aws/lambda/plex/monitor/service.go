package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/cmd/plex"
)

func executePlexCommand(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) error {
	fmt.Println(getTime().String())
	plex.Monitor("false")
	return nil
}

func getTime() time.Time {
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println(err)
	}
	return time.Now().In(location)
}

func main() {
	lambda.Start(executePlexCommand)
}
