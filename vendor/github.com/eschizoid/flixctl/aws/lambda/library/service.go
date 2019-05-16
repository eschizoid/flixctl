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
	fmt.Printf("Exectuing λ with payload: %+v\n", input)
	switch input.Command {
	case "library-jobs":
	case "library-initiate":
	case "library-catalogue":
	case "status":
	case "help":
	}
	fmt.Println("Successfully executed λ library")
}

func main() {
	lambda.Start(executeLibraryCommand)
}
