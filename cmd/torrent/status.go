package torrent

import (
	"fmt"

	"os/exec"

	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	"github.com/eschizoid/flixctl/cmd/plex"
	"github.com/spf13/cobra"
)

var StatusTorrentCmd = &cobra.Command{
	Use:   "status",
	Short: "To Show Torrents Status",
	Long:  `to show the status of the torrents being downloaded`,
	Run: func(cmd *cobra.Command, args []string) {
		status := ec2Service.Status(plex.Session, plex.InstanceID)
		if status == ec2RunningStatus {
			out, err := exec.Command("transmission-remote",
				transmissionHostPort,
				"--authenv",
				"--torrent=all",
				"--list").CombinedOutput()
			if err != nil {
				fmt.Printf("Could not list torrents being downloaded: [%s]\n", err)
			}
			fmt.Println(string(out))
		}
	},
}
