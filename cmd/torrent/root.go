package torrent

import (
	"github.com/spf13/cobra"
)

const (
	transmissionHostPort = "marianoflix.duckdns.org:9091"
	ec2RunningStatus     = "Running"
)

var RootTorrentCmd = &cobra.Command{
	Use:   "torrent",
	Short: "To Control Torrent Client",
}

var keywords string
var quality string

func init() {

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
