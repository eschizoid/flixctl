package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/apex/invoke"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	sess "github.com/aws/aws-sdk-go/aws/session"
	lambdaService "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/eschizoid/flixctl/aws/lambda/models"
	"github.com/go-playground/form"
	"github.com/nlopes/slack"
)

var PlexCommandRegexp = regexp.MustCompile(`^(start|stop|status)$`)

func dispatch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { //nolint:gocyclo
	switch request.HTTPMethod {
	case "POST":
		values, err := url.ParseQuery(request.Body)
		if err != nil {
			return clientError(http.StatusBadRequest)
		}
		slash := new(models.Slash)
		err = form.NewDecoder().Decode(slash, values)
		if err != nil {
			return clientError(http.StatusUnprocessableEntity)
		}
		if VerifySlackRequest(request) {
			return clientError(http.StatusForbidden)
		}
		session := sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		client := lambdaService.New(session, &aws.Config{Region: aws.String("us-east-1")})
		switch slash.Command {
		case "/plex":
			if !PlexCommandRegexp.MatchString(slash.Text) || slash.Text == "" {
				return clientError(http.StatusBadRequest)
			}
			input := models.Input{
				Command:    slash.Token,
				Text:       slash.Text,
				LambdaName: "plex-command-executor",
			}
			invokeLambda(client, input)
		case "/library-jobs", "/library-initiate", "/library-catalogue":
			invokeLambda(client, *new(models.Input))
		case "/movies-search", "/movies-request":
			invokeLambda(client, *new(models.Input))
		case "/shows-search", "/shows-request":
			invokeLambda(client, *new(models.Input))
		case "/torrent-status", "/nzb-status":
			invokeLambda(client, *new(models.Input))
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-type": "application/json"},
			Body:       fmt.Sprintf(`{"response_type": "ephemeral", "text":"Executing %s command"}`, slash.Text),
		}, nil
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func invokeLambda(client *lambdaService.Lambda, input models.Input) {
	err := invoke.InvokeAsync(client, input.LambdaName, input)
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

func VerifySlackRequest(request events.APIGatewayProxyRequest) bool {
	headers := http.Header(request.MultiValueHeaders)
	secretsVerifier, err := slack.NewSecretsVerifier(headers, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(&secretsVerifier, request.Body)
	if err != nil {
		panic(err)
	}
	err = secretsVerifier.Ensure()
	if err != nil {
		panic(err)
	}
	return true
}

func main() {
	lambda.Start(dispatch)
}
