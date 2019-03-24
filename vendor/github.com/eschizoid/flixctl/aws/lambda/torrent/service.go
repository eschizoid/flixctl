package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/plex"
)

func executePlexCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Command {
	case "movies-search":
		plex.Start()
	case "shows-search":
		plex.Stop(os.Getenv("SLACK_NOTIFICATION"))
	case "status":
		plex.Status()
	}
	fmt.Println("Successfully executed plex command")
}

func main() {
	lambda.Start(executePlexCommand)
}
