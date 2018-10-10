package torrent

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	"github.com/eschizoid/flixctl/cmd/plex"
	"github.com/spf13/cobra"
)

var DownloadTorrentCmd = &cobra.Command{
	Use:   "download",
	Short: "To Download a Torrent",
	Long:  `to download a torrent using Transmission client.`,
	Run: func(cmd *cobra.Command, args []string) {
		status := ec2Service.Status(plex.Session, plex.InstanceID)
		magnetLink, err := base64.StdEncoding.DecodeString(os.Getenv("MAGNET_LINK"))
		if err != nil {
			fmt.Printf("Could not decode the magnet link: [%s]\n", err)
		}
		if status == ec2RunningStatus {
			transmission := exec.Command("transmission-remote",
				transmissionHostPort,
				"--authenv",
				"--add",
				string(magnetLink))
			err := transmission.Start()
			if err != nil {
				fmt.Println("Could not download torrent using the given magnet link")
			}
		}
	},
}
