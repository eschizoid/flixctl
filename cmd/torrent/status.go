package torrent

import (
	"github.com/spf13/cobra"
)

var StatusTorrentCmd = &cobra.Command{
	Use:   "status",
	Short: "To Show Torrent Download Status",
	Long:  `to show current download status of all torrents`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
