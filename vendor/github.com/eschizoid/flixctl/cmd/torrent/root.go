package torrent

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	Ec2RunningStatus = "Running"
	Ec2StoppedStatus = "Stopped"
)

var (
	RootTorrentCmd = &cobra.Command{
		Use:   "torrent",
		Short: "To Control Torrent Client",
	}
	slackNotification       string
	magnetLink              string
	keywords                string
	quality                 string
	slackIncomingHookURL    string
	downloadDir             string
	awsResourceTagNameValue = os.Getenv("AWS_RESOURCE_TAG_NAME_VALUE")
)

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
		DownloadTorrentCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
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
		SearchTorrentCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		StatusTorrentCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_TORRENT_INCOMING_HOOK_URL"),
			"slack channel to notify of the torrent event",
		)
		StatusTorrentCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		RootTorrentCmd.AddCommand(SearchTorrentCmd, DownloadTorrentCmd, StatusTorrentCmd)
		return struct{}{}
	}()
)
