package sonarr

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	RootSonarrCmd = &cobra.Command{
		Use:   "sonarr",
		Short: "To Control Sonarr",
	}
	keywords             string
	slackNotification    string
	slackIncomingHookURL string
)

var (
	_ = func() struct{} {
		RootSonarrCmd.AddCommand(SearchSonarrCmd)
		return struct{}{}
	}()
)

var (
	_ = func() struct{} {
		SearchSonarrCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		SearchSonarrCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_REQUESTS_HOOK_URL"),
			"slack channel to notify",
		)
		SearchSonarrCmd.Flags().StringVarP(&keywords,
			"keyword",
			"k",
			"",
			"the keywords that will be used to search for matching title movies",
		)
		RootSonarrCmd.AddCommand(
			SearchSonarrCmd,
		)
		return struct{}{}
	}()
)
