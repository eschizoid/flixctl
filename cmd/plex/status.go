package plex

import (
	"fmt"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	slackService "github.com/eschizoid/flixctl/slack/plex"
	"github.com/spf13/cobra"
)

var StatusPlexCmd = &cobra.Command{
	Use:   "status",
	Short: "To Get Plex Status",
	Long:  "to get the status of the Plex Media Center.",
	Run: func(cmd *cobra.Command, args []string) {
		plexStatus := Status()
		fmt.Println("EC2 current status: " + plexStatus)
	},
}

func Status() string {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	svc := ec2.New(awsSession, awsSession.Config)
	var instanceID = ec2Service.FetchInstanceID(svc, "plex")
	plexStatus := ec2Service.Status(svc, instanceID)
	if slackIncomingHookURL != "" {
		slackService.SendStatus(plexStatus, slackIncomingHookURL)
	}
	return plexStatus
}
