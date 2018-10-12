package plex

import (
	"fmt"

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
		if ec2Service.Status(Session, InstanceID) == "Stopped" {
			slackService.SendStop()
			return
		}
		var oldVolumeID = ebsService.FetchVolumeID(Session, "plex")
		snapService.Create(Session, oldVolumeID, "plex")
		ec2Service.Stop(Session, InstanceID)
		ebsService.Detach(Session, oldVolumeID)
		ebsService.Delete(Session, oldVolumeID)
		slackService.SendStop()
		fmt.Println("Plex Stopped")
	},
}
