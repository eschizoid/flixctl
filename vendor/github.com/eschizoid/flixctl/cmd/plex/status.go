package plex

import (
	sess "github.com/aws/aws-sdk-go/aws/session"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	slackService "github.com/eschizoid/flixctl/slack/plex"
	"github.com/spf13/cobra"
)

var StatusPlexCmd = &cobra.Command{
	Use:   "status",
	Short: "To Get Plex Status",
	Long:  `to get the status of the Plex Media Center.`,
	Run: func(cmd *cobra.Command, args []string) {
		Status()
	},
}

func Status() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	var instanceID = ec2Service.FetchInstanceID(awsSession, "plex")
	plexStatus := ec2Service.Status(awsSession, instanceID)
	slackService.SendStatus(plexStatus)
}
