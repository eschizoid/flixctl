package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/slack"
	"github.com/eschizoid/flixctl/cmd/plex"
	"github.com/eschizoid/flixctl/worker"
	"github.com/go-playground/form"
)

var CommandRegexp = regexp.MustCompile(`^(start|stop|status|token)$`)

func router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return dispatch(request)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func dispatch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { //nolint:gocyclo
	values, err := url.ParseQuery(request.Body)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}
	slash := new(slack.Slash)
	err = form.NewDecoder().Decode(slash, values)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}
	if !CommandRegexp.MatchString(slash.Text) || slash.Text == "" {
		return clientError(http.StatusBadRequest)
	}
	if !slack.VerifySlackRequest(request) {
		return clientError(http.StatusForbidden)
	}

	startTask := func() interface{} {
		plex.Start()
		return "done executing start plex!"
	}
	stopTask := func() interface{} {
		plex.Stop("true")
		return "done executing stop plex!"
	}
	statusTask := func() interface{} {
		plex.Status()
		return "done executing status plex!"
	}
	tokenTask := func() interface{} {
		plex.Token()
		return "done executing token plex!"
	}

	var tasks []worker.TaskFunction
	switch slash.Text {
	case "start":
		tasks = []worker.TaskFunction{startTask}
	case "stop":
		tasks = []worker.TaskFunction{stopTask}
	case "status":
		tasks = []worker.TaskFunction{statusTask}
	case "token":
		tasks = []worker.TaskFunction{tokenTask}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultChannel := worker.PerformTasks(ctx, tasks)

	// Print value from first goroutine and cancel others
	for result := range resultChannel {
		switch v := result.(type) {
		case error:
			fmt.Println(v)
		case string:
			fmt.Println(v)
		default:
			fmt.Println("Some unknown type ")
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       fmt.Sprintf(`{"response_type": "ephemeral", "text":"Executing %s command"}`, slash.Text),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
