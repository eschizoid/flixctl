package torrent

import (
	"github.com/spf13/cobra"
)

var StatusTorrentCmd = &cobra.Command{
	Use:   "status",
	Short: "To Show the Status of the Torrents",
	Long:  `to show the status of the torrents being downloaded`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
