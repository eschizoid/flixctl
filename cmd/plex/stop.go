package plex

import (
	"fmt"
	"os"

	ebsService "github.com/eschizoid/flixctl/aws/ebs"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	snapService "github.com/eschizoid/flixctl/aws/snapshot"
	"github.com/spf13/cobra"
)

var StopPlexCmd = &cobra.Command{
	Use:   "stop",
	Short: "To Stop Plex",
	Long:  `to stop the Plex Media Center.`,
	Run: func(cmd *cobra.Command, args []string) {
		if ec2Service.Status(Session, InstanceID) == "Stopped" {
			os.Exit(0)
		}
		var oldVolumeID = ebsService.FetchVolumeID(Session, "plex")
		snapService.Create(Session, oldVolumeID, "plex")
		ec2Service.Stop(Session, InstanceID)
		ebsService.Detach(Session, oldVolumeID)
		ebsService.Delete(Session, oldVolumeID)
		fmt.Println("Plex Stopped")
	},
}
