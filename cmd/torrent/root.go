package torrent

import (
	"github.com/spf13/cobra"
)

var RootTorrentCmd = &cobra.Command{
	Use:   "torrent",
	Short: "To Control Torrent Client",
}

var argMagnetLink string
var keywords string
var quality string

func init() {
	DownloadTorrentCmd.Flags().StringVarP(&argMagnetLink,
		"magnet-link",
		"m",
		"",
		"uri of the torrent magnet link to download",
	)
	SearchTorrentCmd.Flags().StringVarP(&keywords,
		"keywords",
		"k",
		"",
		"the keywords that will be used to search for available torrents",
	)
	SearchTorrentCmd.Flags().StringVarP(&quality,
		"minimum-quality",
		"q",
		"",
		"the minimum quality of the torrent file",
	)
	RootTorrentCmd.AddCommand(SearchTorrentCmd, DownloadTorrentCmd, StatusTorrentCmd)
}
