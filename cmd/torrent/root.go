package torrent

import (
	"github.com/spf13/cobra"
)

var RootTorrentCmd = &cobra.Command{
	Use:   "torrent",
	Short: "To Control Torrent Client",
}

var magnetLink string
var keywords string

func init() {

	DownloadTorrentCmd.Flags().StringVarP(&magnetLink,
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

	RootTorrentCmd.AddCommand(SearchTorrentCmd, DownloadTorrentCmd, StatusTorrentCmd)
}
