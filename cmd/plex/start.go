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

var StartPlexCmd = &cobra.Command{
	Use:   "start",
	Short: "To Start Plex",
	Long:  `to start the Plex Media Center.`,
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

func Start() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	var instanceID = ec2Service.FetchInstanceID(awsSession, "plex")
	if ec2Service.Status(awsSession, instanceID) == "Running" {
		slackService.SendStart()
		return
	}
	ec2Service.Start(awsSession, instanceID)
	var oldSnapshotID = snapService.FetchSnapshotID(awsSession, "plex")
	ebsService.Create(awsSession, oldSnapshotID, "plex")
	var newVolumeID = ebsService.FetchVolumeID(awsSession, "plex")
	ebsService.Attach(awsSession, instanceID, newVolumeID)
	snapService.Delete(awsSession, oldSnapshotID)
	slackService.SendStart()
	fmt.Println("Plex Running")
}
