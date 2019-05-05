package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
)

func executeMoviesCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	fmt.Printf("Exectuing λ with payload: %+v\n", input)
	switch input.Command {
	case "shows-search":
	case "shows-request":
	}
	fmt.Println("Successfully executed plex command")
}

func main() {
	lambda.Start(executeMoviesCommand)
}
