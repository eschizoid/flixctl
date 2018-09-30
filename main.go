package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	ebsService "github.com/eschizoid/flixctl/aws/ebs"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	snapService "github.com/eschizoid/flixctl/aws/snap"
	"github.com/spf13/cobra"
	"gopkg.in/cheggaaa/pb.v1"
	"time"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var instanceId = ec2Service.FetchInstanceId(sess, "plex")

	var startPlexCmd = &cobra.Command{
		Use:   "start",
		Short: "To Start Plex",
		Long:  `to start the Plex Media Center.`,
		Run: func(cmd *cobra.Command, args []string) {
			ec2Service.Start(sess, instanceId)
			//TODO Poll instead of Wait
			println("Starting for EC2 instance")
			wait()
			var oldSnapshotId = snapService.FetchSnapshotId(sess, "plex")
			ebsService.Create(sess, oldSnapshotId)
			//TODO Poll instead of Wait
			println("Creating for EBS volume")
			wait()
			var newVolumeId = ebsService.FetchVolumeId(sess, "plex")
			ebsService.Attach(sess, instanceId, newVolumeId)
			//TODO Poll instead of Wait
			println("Attaching for EBS volume")
			wait()
			snapService.Delete(sess, oldSnapshotId)
		},
	}

	var stopPlexCmd = &cobra.Command{
		Use:   "stop",
		Short: "To Stop Plex",
		Long:  `to stop the Plex Media Center.`,
		Run: func(cmd *cobra.Command, args []string) {
			var oldVolumeId = ebsService.FetchVolumeId(sess, "plex")
			snapService.Create(sess, oldVolumeId)
			//TODO Poll instead of Wait
			println("Creating Snapshot")
			wait()
			ec2Service.Stop(sess, instanceId)
			ebsService.Detach(sess, oldVolumeId)
			//TODO Poll instead of Wait
			println("Detaching EBS volume")
			wait()
			ebsService.Delete(sess, oldVolumeId)
		},
	}

	var statusPlexCmd = &cobra.Command{
		Use:   "status",
		Short: "To Get Plex Status",
		Long:  `to get the status of the Plex Media Center.`,
		Run: func(cmd *cobra.Command, args []string) {
			ec2Service.Status(sess, instanceId)
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
	flixctlCmd.Execute()
}

func wait() {
	count := 30
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Second)
	}
	bar.FinishPrint("Done!")
}
