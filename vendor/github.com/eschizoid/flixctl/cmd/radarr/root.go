package radarr

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	RootRadarrCmd = &cobra.Command{
		Use:   "radarr",
		Short: "To Control Radarr",
	}
	keywords             string
	slackNotification    string
	slackIncomingHookURL string
)

var (
	_ = func() struct{} {
		RootRadarrCmd.AddCommand(SearchRadarrCmd)
		return struct{}{}
	}()
)

var (
	_ = func() struct{} {
		SearchRadarrCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		SearchRadarrCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_REQUESTS_HOOK_URL"),
			"slack channel to notify",
		)
		SearchRadarrCmd.Flags().StringVarP(&keywords,
			"keyword",
			"k",
			"",
			"the keywords that will be used to search for matching title shows",
		)
		RootRadarrCmd.AddCommand(
			SearchRadarrCmd,
		)
		return struct{}{}
	}()
)
