package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	ebsService "github.com/eschizoid/flixctl/aws/ebs"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	snapService "github.com/eschizoid/flixctl/aws/snap"
	"github.com/spf13/cobra"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var instanceID = ec2Service.FetchInstanceID(sess, "plex")

	var startPlexCmd = &cobra.Command{
		Use:   "start",
		Short: "To Start Plex",
		Long:  `to start the Plex Media Center.`,
		Run: func(cmd *cobra.Command, args []string) {
			if ec2Service.Status(sess, instanceID) == "Running" {
				os.Exit(0)
			}
			ec2Service.Start(sess, instanceID)
			var oldSnapshotID = snapService.FetchSnapshotID(sess, "plex")
			ebsService.Create(sess, oldSnapshotID, "plex")
			var newVolumeID = ebsService.FetchVolumeID(sess, "plex")
			ebsService.Attach(sess, instanceID, newVolumeID)
			snapService.Delete(sess, oldSnapshotID)
			fmt.Println("Plex Running")
		},
	}

	var stopPlexCmd = &cobra.Command{
		Use:   "stop",
		Short: "To Stop Plex",
		Long:  `to stop the Plex Media Center.`,
		Run: func(cmd *cobra.Command, args []string) {
			if ec2Service.Status(sess, instanceID) == "Stopped" {
				os.Exit(0)
			}
			var oldVolumeID = ebsService.FetchVolumeID(sess, "plex")
			snapService.Create(sess, oldVolumeID, "plex")
			ec2Service.Stop(sess, instanceID)
			ebsService.Detach(sess, oldVolumeID)
			ebsService.Delete(sess, oldVolumeID)
			fmt.Println("Plex Stopped")
		},
	}

	var statusPlexCmd = &cobra.Command{
		Use:   "status",
		Short: "To Get Plex Status",
		Long:  `to get the status of the Plex Media Center.`,
		Run: func(cmd *cobra.Command, args []string) {
			ec2Service.Status(sess, instanceID)
		},
	}

	var downloadTorrentCmd = &cobra.Command{
		Use:   "download",
		Short: "To Download a Torrent",
		Long:  `to download a torrent using Transmission client.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var statusTorrentCmd = &cobra.Command{
		Use:   "status",
		Short: "To Get Current Torrent Downloads",
		Long:  `to get all Transmission current downloads.`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var flixctlCmd = &cobra.Command{Use: "flixctl"}
	var plexCmd = &cobra.Command{Use: "plex", Short: "To Control Plex Media Center"}
	var torrentCmd = &cobra.Command{Use: "torrent", Short: "To Control Torrent Client"}

	plexCmd.AddCommand(startPlexCmd, stopPlexCmd, statusPlexCmd)
	flixctlCmd.AddCommand(plexCmd)
	torrentCmd.AddCommand(downloadTorrentCmd, statusTorrentCmd)
	flixctlCmd.AddCommand(torrentCmd)
	if err := flixctlCmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
