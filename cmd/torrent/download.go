package torrent

import (
	"github.com/spf13/cobra"
)

var DownloadTorrentCmd = &cobra.Command{
	Use:   "download",
	Short: "To Download a Torrent",
	Long:  `to download a torrent using Transmission client.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
