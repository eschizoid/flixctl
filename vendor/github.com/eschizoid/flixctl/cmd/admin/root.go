package admin

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	FlixctlAdminRootCmd = &cobra.Command{
		Use:   "admin",
		Short: "To Perform Admin and Maintenance Tasks",
		Long:  "to perform admin / maintenance tasks",
	}
	slackClientID     string
	slackClientSecret string
	slackCode         string
	slackRedirectURI  string
)

var (
	_ = func() struct{} {
		SlackOauthTokenCmd.Flags().StringVarP(&slackClientID,
			"slack-client-id",
			"i",
			os.Getenv("SLACK_CLIENT_ID"),
			"slack client id",
		)
		SlackOauthTokenCmd.Flags().StringVarP(&slackClientSecret,
			"slack-client-secret",
			"s",
			os.Getenv("SLACK_CLIENT_SECRET"),
			"slack client secret",
		)
		SlackOauthTokenCmd.Flags().StringVarP(&slackCode,
			"slack-code",
			"c",
			os.Getenv("SLACK_CODE"),
			"slack code",
		)
		SlackOauthTokenCmd.Flags().StringVarP(&slackRedirectURI,
			"slack-redirect-uri",
			"r",
			os.Getenv("SLACK_REDIRECT_URI"),
			"slack redirect uri",
		)
		FlixctlAdminRootCmd.AddCommand(
			PurgeSlackCmd,
			RenewCertsCmd,
			RestartPlexServicesCmd,
			SlackOauthTokenCmd,
		)
		return struct{}{}
	}()
)

func ShowError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
