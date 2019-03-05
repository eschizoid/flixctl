package library

import (
	"encoding/json"
	"fmt"

	slackService "github.com/eschizoid/flixctl/slack/auth"
	"github.com/nlopes/slack"
	"github.com/spf13/cobra"
)

var TokenCmd = &cobra.Command{
	Use:   "token",
	Short: "To Get An Oauth Token",
	Long:  "to get an oauth token for a given client id",
	Run: func(cmd *cobra.Command, args []string) {
		if resp, err := slack.GetOAuthResponse(slackClientID, slackClientSecret, slackCode, slackRedirectURI, true); err != nil {
			m := make(map[string]string)
			m["token_saved_successfully"] = "false"
			jsonString, _ := json.Marshal(m)
			fmt.Println("\n" + string(jsonString))
		} else {
			if err := slackService.SaveToken(slackClientID, resp.AccessToken); err != nil {
				ShowError(err)
			}
			m := make(map[string]string)
			m["token_saved_successfully"] = "true"
			jsonString, _ := json.Marshal(m)
			fmt.Println("\n" + string(jsonString))
		}
	},
}
