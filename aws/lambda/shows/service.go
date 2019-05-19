//nolint:dupl
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/ombi"
	"github.com/eschizoid/flixctl/cmd/sonarr"
	slackService "github.com/eschizoid/flixctl/slack/sonarr"
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
			slackService.SendShowsHelp(os.Getenv("SLACK_GENERAL_HOOK_URL"))
		} else {
			sonarr.SearchShows()
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "shows-request":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" {
			slackService.SendShowsHelp(os.Getenv("SLACK_GENERAL_HOOK_URL"))
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
