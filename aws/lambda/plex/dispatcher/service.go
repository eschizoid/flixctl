package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/apex/invoke"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	sess "github.com/aws/aws-sdk-go/aws/session"
	lambdaService "github.com/aws/aws-sdk-go/service/lambda"
	types "github.com/eschizoid/flixctl/aws/lambda"
	"github.com/go-playground/form"
)

var CommandRegexp = regexp.MustCompile(`^(start|stop|status)$`)

func router(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "POST":
		return dispatch(request)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func dispatch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	values, err := url.ParseQuery(request.Body)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}
	slash := new(types.Slash)
	err = form.NewDecoder().Decode(slash, values)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}
	if !CommandRegexp.MatchString(slash.Text) || slash.Text == "" {
		return clientError(http.StatusBadRequest)
	}
	invokeLambda(slash)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       fmt.Sprintf(`{"response_type": "in_channel", "text":"Executing %s command"}`, slash.Text),
	}, nil
}

func invokeLambda(slash *types.Slash) {
	session := sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	client := lambdaService.New(session, &aws.Config{Region: aws.String("us-east-1")})
	switch slash.Text {
	case "start":
		err := invoke.InvokeAsync(client, "plex-command-executor", types.Input{Command: "start"})
		if err != nil {
			fmt.Println("Error invoking λ:", err)
		}
	case "stop":
		err := invoke.InvokeAsync(client, "plex-command-executor", types.Input{Command: "stop"})
		if err != nil {
			fmt.Println("Error invoking λ:", err)
		}
	case "status":
		err := invoke.InvokeAsync(client, "plex-command-executor", types.Input{Command: "status"})
		if err != nil {
			fmt.Println("Error invoking λ:", err)
		}
	}
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
