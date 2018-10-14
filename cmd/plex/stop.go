package plex

import (
	"fmt"

	sess "github.com/aws/aws-sdk-go/aws/session"
	ebsService "github.com/eschizoid/flixctl/aws/ebs"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	snapService "github.com/eschizoid/flixctl/aws/snapshot"
	slackService "github.com/eschizoid/flixctl/slack/plex"
	"github.com/spf13/cobra"
)

var StopPlexCmd = &cobra.Command{
	Use:   "stop",
	Short: "To Stop Plex",
	Long:  `to stop the Plex Media Center.`,
	Run: func(cmd *cobra.Command, args []string) {
		Stop()
	},
}

func Stop() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	var instanceID = ec2Service.FetchInstanceID(awsSession, "plex")
	if ec2Service.Status(awsSession, instanceID) == "Stopped" {
		slackService.SendStop()
		return
	}
	var oldVolumeID = ebsService.FetchVolumeID(awsSession, "plex")
	snapService.Create(awsSession, oldVolumeID, "plex")
	ec2Service.Stop(awsSession, instanceID)
	ebsService.Detach(awsSession, oldVolumeID)
	ebsService.Delete(awsSession, oldVolumeID)
	slackService.SendStop()
	fmt.Println("Plex Stopped")
}
