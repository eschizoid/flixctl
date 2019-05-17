package main

import (
	"encoding/json"
	"flixctl/cmd/ombi"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/sonarr"
	slackLambdaService "github.com/eschizoid/flixctl/slack/lambda"
)

func executeShowsCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	fmt.Printf("Exectuing λ with payload: %+v\n", input)
	switch input.Command {
	case "shows-search":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" { //nolint:goconst
			slackLambdaService.SendShowsHelp("")
		} else {
			sonarr.SearchShows()
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "shows-request":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" {
			slackLambdaService.SendShowsHelp("")
		} else {
			ombi.Request()
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	}
	fmt.Println("Successfully executed λ shows")
}

func main() {
	lambda.Start(executeShowsCommand)
}
