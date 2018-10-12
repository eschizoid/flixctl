package torrent

import (
	"github.com/eschizoid/flixctl/slack"
	"github.com/eschizoid/flixctl/torrent"
	"github.com/spf13/cobra"
)

var StatusTorrentCmd = &cobra.Command{
	Use:   "status",
	Short: "To Show Torrents Status",
	Long:  `to show the status of the torrents being downloaded`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadStatus := torrent.Status()
		slack.SendStatus(downloadStatus)
	},
}
