package torrent

import (
	"fmt"

	slackService "github.com/eschizoid/flixctl/slack/torrent"
	"github.com/eschizoid/flixctl/torrent"
	"github.com/spf13/cobra"
)

var StatusTorrentCmd = &cobra.Command{
	Use:   "status",
	Short: "To Show Torrents Status",
	Long:  `to show the status of the torrents being downloaded`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadStatus := torrent.Status()
		if slackIncomingHookURL != "" {
			slackService.SendStatus(downloadStatus, slackIncomingHookURL)
		}
		fmt.Println(downloadStatus)
	},
}
