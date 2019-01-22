package slack

import (
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
)

func VerifySlackRequest(request events.APIGatewayProxyRequest) bool {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	headers := http.Header(request.MultiValueHeaders)
	secretsVerifier, err := slack.NewSecretsVerifier(headers, signingSecret)
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
