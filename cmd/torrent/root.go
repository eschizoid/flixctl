package torrent

import (
	"github.com/spf13/cobra"
)

var RootTorrentCmd = &cobra.Command{
	Use:   "torrent",
	Short: "To Control Torrent Client",
}

func init() {
	RootTorrentCmd.AddCommand(DownloadTorrentCmd, StatusTorrentCmd)
}
