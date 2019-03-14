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
	//maxInactiveTime         string
	slackNotification       string
	slackIncomingHookURL    string
	awsResourceTagNameValue = os.Getenv("AWS_RESOURCE_TAG_NAME_VALUE")
)

var (
	_ = func() struct{} {
		StartPlexCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_MONITORING_HOOK_URL"),
			"slack channel to notify",
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
			os.Getenv("SLACK_MONITORING_HOOK_URL"),
			"slack channel to notify",
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
			os.Getenv("SLACK_MONITORING_HOOK_URL"),
			"slack channel to notify",
		)
		StatusPlexCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		RootPlexCmd.AddCommand(
			StartPlexCmd,
			StatusPlexCmd,
			StopPlexCmd,
			TokenPlexCmd,
		)
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

func ShowError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
