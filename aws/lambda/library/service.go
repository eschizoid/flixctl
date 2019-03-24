package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
)

func executeLibraryCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	switch input.Command {
	case "library-jobs":
	case "library-initiate":
	case "library-catalogue":
	case "status":
	}
	fmt.Println("Successfully executed plex command")
}

func main() {
	lambda.Start(executeLibraryCommand)
}
