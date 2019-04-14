package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

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

func dispatch(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) { //nolint:gocyclo

	switch request.HTTPMethod {
	case "POST":
		if !isValidSlackRequest(request) {
			return clientError(http.StatusForbidden)
		}
	default:
		return clientError(http.StatusMethodNotAllowed)
	}

	session := sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	client := lambdaService.New(session, &aws.Config{Region: aws.String("us-east-1")})

	values, err := url.ParseQuery(request.Body)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}

	slash := new(models.Slash)
	err = form.NewDecoder().Decode(slash, values)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	lambdaName := request.QueryStringParameters["lambda"]
	command := strings.Replace(slash.Command, "/", "", -1)
	text := slash.Text
	input := models.Input{
		Command:    command,
		Parameter:  text,
		LambdaName: lambdaName,
	}
	fmt.Printf("Invoking λ with payload: %+v\n", lambdaName)

	invokeLambda(client, lambdaName, input)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-type": "application/json"},
		Body:       fmt.Sprintf(`{"response_type": "ephemeral", "text":"Executing command [%s] with parameters []"}`, command),
	}, nil
}

func invokeLambda(client *lambdaService.Lambda, lambdaName string, input interface{}) {
	if err := invoke.InvokeAsyncQualifier(client, lambdaName, "$LATEST", input); err != nil {
		fmt.Println("Error invoking λ: ", err)
	}
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func isValidSlackRequest(request events.APIGatewayProxyRequest) bool {
	headers := http.Header(request.MultiValueHeaders)
	secretsVerifier, err := slack.NewSecretsVerifier(headers, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		fmt.Println("Error invoking λ: ", err)
		return false
	}
	if _, err = io.WriteString(&secretsVerifier, request.Body); err != nil {
		fmt.Println("Error invoking λ: ", err)
		return false
	}
	if err := secretsVerifier.Ensure(); err != nil {
		fmt.Println("Error invoking λ: ", err)
		return false
	}
	return true
}

func main() {
	lambda.Start(dispatch)
}
