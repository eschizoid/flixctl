package plex

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	Ec2RunnningStatus = "Running"
	Ec2StoppedStatus  = "Stopped"
)

var RootPlexCmd = &cobra.Command{
	Use:   "plex",
	Short: "To Control Plex Media Center",
}

var slackIncomingHookURL string

var (
	_ = func() struct{} {
		StartPlexCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_PLEX_INCOMING_HOOK_URL"),
			"slack channel to notify of the plex event",
		)
		StopPlexCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_PLEX_INCOMING_HOOK_URL"),
			"slack channel to notify of the plex event",
		)
		StatusPlexCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_PLEX_INCOMING_HOOK_URL"),
			"slack channel to notify of the plex event",
		)
		RootPlexCmd.AddCommand(StartPlexCmd, StopPlexCmd, StatusPlexCmd)
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
