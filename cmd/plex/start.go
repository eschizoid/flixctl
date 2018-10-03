package plex

import (
	"fmt"
	"os"

	ebsService "github.com/eschizoid/flixctl/aws/ebs"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	snapService "github.com/eschizoid/flixctl/aws/snapshot"
	"github.com/spf13/cobra"
)

var StartPlexCmd = &cobra.Command{
	Use:   "start",
	Short: "To Start Plex",
	Long:  `to start the Plex Media Center.`,
	Run: func(cmd *cobra.Command, args []string) {
		if ec2Service.Status(Session, InstanceID) == "Running" {
			os.Exit(0)
		}
		ec2Service.Start(Session, InstanceID)
		var oldSnapshotID = snapService.FetchSnapshotID(Session, "plex")
		ebsService.Create(Session, oldSnapshotID, "plex")
		var newVolumeID = ebsService.FetchVolumeID(Session, "plex")
		ebsService.Attach(Session, InstanceID, newVolumeID)
		snapService.Delete(Session, oldSnapshotID)
		fmt.Println("Plex Running")
	},
}
