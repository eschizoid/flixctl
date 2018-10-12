package plex

import (
	"fmt"

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
		if ec2Service.Status(Session, InstanceID) == "Running" {
			slackService.SendStart()
			return
		}
		ec2Service.Start(Session, InstanceID)
		var oldSnapshotID = snapService.FetchSnapshotID(Session, "plex")
		ebsService.Create(Session, oldSnapshotID, "plex")
		var newVolumeID = ebsService.FetchVolumeID(Session, "plex")
		ebsService.Attach(Session, InstanceID, newVolumeID)
		snapService.Delete(Session, oldSnapshotID)
		slackService.SendStart()
		fmt.Println("Plex Running")
	},
}
