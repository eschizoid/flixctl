//nolint:dupl
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/eschizoid/flixctl/cmd/ombi"
	"github.com/eschizoid/flixctl/cmd/radarr"
	slackService "github.com/eschizoid/flixctl/slack/radarr"
)

func executeMoviesCommand(evt json.RawMessage) {
	var input models.Input
	if err := json.Unmarshal(evt, &input); err != nil {
		panic(err)
	}
	fmt.Printf("Exectuing λ with payload: %+v\n", input)
	switch input.Command {
	case "movies-search":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" { //nolint:goconst
			slackService.SendMoviesHelp(os.Getenv("SLACK_GENERAL_HOOK_URL"))
		} else {
			radarr.SearchMovies()
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	case "movies-request":
		fmt.Printf("Executing %s command \n", input.Argument)
		if input.Argument == "help" {
			slackService.SendMoviesHelp(os.Getenv("SLACK_GENERAL_HOOK_URL"))
		} else {
			ombi.Request()
		}
		fmt.Printf("Succesfully executed %s \n", input.Argument)
	}
	fmt.Println("Successfully executed λ movies")
}

func main() {
	lambda.Start(executeMoviesCommand)
}
