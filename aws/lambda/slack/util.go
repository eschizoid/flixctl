package slack

import (
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
)

var SigningSecret = os.Getenv("SLACK_SIGNING_SECRET")

func VerifySlackRequest(request events.APIGatewayProxyRequest) bool {
	headers := http.Header(request.MultiValueHeaders)
	secretsVerifier, err := slack.NewSecretsVerifier(headers, SigningSecret)
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
