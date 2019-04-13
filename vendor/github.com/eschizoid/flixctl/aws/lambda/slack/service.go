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
		if VerifySlackRequest(request) {
			return clientError(http.StatusForbidden)
		}
	default:
		return clientError(http.StatusMethodNotAllowed)
	}

	session := sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	client := lambdaService.New(session, &aws.Config{Region: aws.String("us-east-1")})

	if _, fromSlackButton := request.QueryStringParameters["lambda-name"]; fromSlackButton {
		values, err := url.ParseQuery(request.Body)
		if err != nil {
			return clientError(http.StatusBadRequest)
		}
		slash := new(models.Slash)
		err = form.NewDecoder().Decode(slash, values)
		if err != nil {
			return clientError(http.StatusUnprocessableEntity)
		}
		switch slash.Command {
		case "/plex":
			if !PlexCommandRegexp.MatchString(slash.Text) || slash.Text == "" {
				return clientError(http.StatusBadRequest)
			}
			input := models.Input{
				Command:    slash.Token,
				Text:       slash.Text,
				LambdaName: "plex-executor",
			}
			invokeLambda(client, "plex-executor", input)
		case "/library-jobs", "/library-initiate", "/library-catalogue":
			input := models.Input{
				Command:    slash.Token,
				Text:       slash.Text,
				LambdaName: "library-executor",
			}
			invokeLambda(client, "library-executor", input)
		case "/movies-search", "/shows-search":
			input := models.Input{
				Command:    slash.Token,
				Text:       slash.Text,
				LambdaName: "torrent-search-executor",
			}
			invokeLambda(client, "torrent-search-executor", input)
		case "/torrent-status", "/nzb-status":
			input := models.Input{
				Command: slash.Token,
				Text:    slash.Text,
			}
			invokeLambda(client, "torrent-status-executor", input)
		}
	} else {
		switch request.QueryStringParameters["lambda-name"] {
		case "torrent-download":
			invokeLambda(client, "torrent-download-executor", request.QueryStringParameters)
		case "torrent-request":
			invokeLambda(client, "torrent-download-executor", request.QueryStringParameters)
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       fmt.Sprintf(`{"response_type": "ephemeral", "text":"Executing command"}`),
	}, nil
}

func invokeLambda(client *lambdaService.Lambda, lambdaName string, input interface{}) {
	err := invoke.InvokeAsync(client, lambdaName, input)
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
