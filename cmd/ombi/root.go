package ombi

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	RootOmbiCmd = &cobra.Command{
		Use:   "ombi",
		Short: "To Control Ombi",
		Long:  "to control ombi",
	}
	keywords             string
	slackNotification    string
	slackIncomingHookURL string
)

var (
	_ = func() struct{} {
		RequestOmbiCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		RequestOmbiCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_REQUESTS_HOOK_URL"),
			"slack channel to notify",
		)
		RequestOmbiCmd.Flags().StringVarP(&keywords,
			"keyword",
			"k",
			"",
			"the keywords that will be used to search for matching title movies / shows",
		)
		RootOmbiCmd.AddCommand(
			RequestOmbiCmd,
		)
		return struct{}{}
	}()
)
