package library

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	slackService "github.com/eschizoid/flixctl/slack/auth"
	"github.com/nlopes/slack"
	"github.com/spf13/cobra"
)

var TokenCmd = &cobra.Command{
	Use:   "token",
	Short: "To Get An Oauth Token",
	Long:  "to get an oauth token for a given client id",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		awsSession.Config.Endpoint = aws.String("http://dynamodb:8000")
		svc := dynamodb.New(awsSession)
		if resp, err := slack.GetOAuthResponse(slackClientID, slackClientSecret, slackCode, slackRedirectURI, true); err != nil {
			m := make(map[string]string)
			m["slack_oauth_token_saved"] = "false"
			jsonString, _ := json.Marshal(m)
			fmt.Println("\n" + string(jsonString))
		} else {
			if err := slackService.SaveToken(slackClientID, resp.AccessToken, svc); err != nil {
				ShowError(err)
			}
			m := make(map[string]string)
			m["slack_oauth_token_saved"] = "true"
			jsonString, _ := json.Marshal(m)
			fmt.Println("\n" + string(jsonString))
		}
	},
}
