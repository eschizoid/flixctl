package plex

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	Ec2RunningStatus = "Running"
	Ec2StoppedStatus = "Stopped"
)

var (
	RootPlexCmd = &cobra.Command{
		Use:   "plex",
		Short: "To Control Plex Media Center",
	}
	slackNotification    string
	slackIncomingHookURL string
)

var (
	_ = func() struct{} {
		StartPlexCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_PLEX_INCOMING_HOOK_URL"),
			"slack channel to notify of the plex event",
		)
		StartPlexCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		StopPlexCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_PLEX_INCOMING_HOOK_URL"),
			"slack channel to notify of the plex event",
		)
		StopPlexCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		StatusPlexCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_PLEX_INCOMING_HOOK_URL"),
			"slack channel to notify of the plex event",
		)
		StatusPlexCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		RootPlexCmd.AddCommand(StartPlexCmd, StopPlexCmd, StatusPlexCmd, TokenPlexCmd)
		return struct{}{}
	}()
)

func Indicator(shutdownCh <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case <-shutdownCh:
			return
		}
	}
}
