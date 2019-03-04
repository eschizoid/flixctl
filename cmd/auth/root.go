package library

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	OauthSlackRootCmd = &cobra.Command{
		Use:   "auth",
		Short: "To Start Oauth Flow With Slack",
	}
	slackClientID     string
	slackClientSecret string
	slackCode         string
	slackRedirectURI  string
)

var (
	_ = func() struct{} {
		TokenCmd.Flags().StringVarP(&slackClientID,
			"slackClientId",
			"c",
			os.Getenv("SLACK_CLIENT_ID"),
			"slack client id",
		)
		TokenCmd.Flags().StringVarP(&slackClientSecret,
			"slack-client-secret",
			"s",
			os.Getenv("SLACK_CLIENT_SECRET"),
			"slack client secret",
		)
		TokenCmd.Flags().StringVarP(&slackCode,
			"slack-code",
			"o",
			os.Getenv("SLACK_CODE"),
			"slack code",
		)
		TokenCmd.Flags().StringVarP(&slackRedirectURI,
			"slack-redirect-uri",
			"r",
			os.Getenv("SLACK_REDIRECT_URI"),
			"slack redirect uri",
		)
		OauthSlackRootCmd.AddCommand(TokenCmd)
		return struct{}{}
	}()
)

func ShowError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
