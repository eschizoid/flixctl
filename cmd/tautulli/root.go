package tautulli

import (
	"os"

	"github.com/spf13/cobra"
)

var RootTautulliCmd = &cobra.Command{
	Use:   "tautulli",
	Short: "Used Internally To Forward Tautulli Notifications",
}

var message string
var slackIncomingHookURL string

func init() {
	ForwardCmd.Flags().StringVarP(&message,
		"message",
		"m",
		"",
		"tautulli message to forward",
	)
	ForwardCmd.Flags().StringVarP(&slackIncomingHookURL,
		"slack-notification-channel",
		"s",
		os.Getenv("SLACK_TAUTULLI_INCOMING_HOOK_URL"),
		"slack channel to notify of the tautulli event",
	)
	RootTautulliCmd.AddCommand(ForwardCmd)
}
