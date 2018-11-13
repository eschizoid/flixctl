package torrent

import (
	"os"

	"github.com/spf13/cobra"
)

var RootTorrentCmd = &cobra.Command{
	Use:   "torrent",
	Short: "To Control Torrent Client",
}

var magnetLink string
var keywords string
var quality string
var slackIncomingHookURL string
var downloadDir string

var (
	_ = func() struct{} {
		DownloadTorrentCmd.Flags().StringVarP(&downloadDir,
			"download-dir",
			"w",
			os.Getenv("DOWNLOAD_DIR"),
			"set the torrent's download folder",
		)
		DownloadTorrentCmd.Flags().StringVarP(&magnetLink,
			"magnet-link",
			"m",
			"",
			"uri of the torrent magnet link to download",
		)
		DownloadTorrentCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_TORRENT_INCOMING_HOOK_URL"),
			"slack channel to notify of the torrent event",
		)
		SearchTorrentCmd.Flags().StringVarP(&downloadDir,
			"download-dir",
			"w",
			os.Getenv("DOWNLOAD_DIR"),
			"set the torrent's download folder",
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
		SearchTorrentCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_TORRENT_INCOMING_HOOK_URL"),
			"slack channel to notify of the torrent event",
		)
		StatusTorrentCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_TORRENT_INCOMING_HOOK_URL"),
			"slack channel to notify of the torrent event",
		)
		RootTorrentCmd.AddCommand(SearchTorrentCmd, DownloadTorrentCmd, StatusTorrentCmd)
		return struct{}{}
	}()
)
