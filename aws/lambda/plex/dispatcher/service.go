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
	"github.com/eschizoid/flixctl/aws/lambda/slack"
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
	session := sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	client := lambdaService.New(session, &aws.Config{Region: aws.String("us-east-1")})
	switch slash.Text {
	case "start":
		invokeLambda(client, "start")
	case "stop":
		invokeLambda(client, "stop")
	case "status":
		invokeLambda(client, "status")
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       fmt.Sprintf(`{"response_type": "ephemeral", "text":"Executing %s command"}`, slash.Text),
	}, nil
}

func invokeLambda(client *lambdaService.Lambda, operation string) {
	err := invoke.InvokeAsync(client, "plex-command-executor", slack.Input{Command: operation})
	if err != nil {
		fmt.Println("Error invoking λ:", err)
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
