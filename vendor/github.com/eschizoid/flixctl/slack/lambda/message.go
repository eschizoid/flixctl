package lambda

import (
	"encoding/json"
	"fmt"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/nlopes/slack"
)

func SendAdminHelp(slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, 4)
	attachmentRenewCerts := slack.AttachmentField{
		Value: `✅ Renew tls certificates all plex related services:\n` + "`/admin renew-certs`",
		Short: false,
	}
	attachmentRestartServices := slack.AttachmentField{
		Value: `✅ Restart all plex related services:\n` + "`/admin restart-services`",
		Short: false,
	}
	attachmentSlackToken := slack.AttachmentField{
		Value: `✅ Get a slack oauth token for a given client id:\n` + "`/admin slack-token`",
		Short: false,
	}
	attachmentSlackPurge := slack.AttachmentField{
		Value: `✅ Purge slack messages from all channels:\n` + "`/admin slack-purge`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color:     "#C40203",
		Title:     "👋 Need some help with" + "`/admin` 💥?",
		TitleLink: util.TorrentDownloadHookURL,
		Fields: []slack.AttachmentField{
			attachmentRenewCerts,
			attachmentRestartServices,
			attachmentSlackToken,
			attachmentSlackPurge,
		},
	}
	attachments = append(attachments, attachment)

	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending download links: [%s]\n", err)
	}
}

func SendMoviesHelp(slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, 2)
	attachmentRequestMovie := slack.AttachmentField{
		Value: `✅ Request movies via ombi:\n` + "`/movies-request`",
		Short: false,
	}
	attachmentSearchMovie := slack.AttachmentField{
		Value: `✅ Search movies using radarr client:\n` + "`/movies-search`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color:     "#C40203",
		Title:     "👋 Need some help with" + "`/movies-request` or `/movies-search`?",
		TitleLink: util.TorrentDownloadHookURL,
		Fields: []slack.AttachmentField{
			attachmentRequestMovie,
			attachmentSearchMovie,
		},
	}
	attachments = append(attachments, attachment)
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending download links: [%s]\n", err)
	}
}


func SendPlexHelp(slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, 4)
	attachmentPlexDisableMonitoring := slack.AttachmentField{
		Value: `✅ Enable Plex monitoring:\n` + "`/plex disable-monitoring`",
		Short: false,
	}
	attachmentPlexEnableMonitoring := slack.AttachmentField{
		Value: `✅ Disable Plex monitoring:\n` + "`/plex enable-monitoring`",
		Short: false,
	}
	attachmentPlexStart := slack.AttachmentField{
		Value: `✅ Start Plex:\n` + "`/plex start`",
		Short: false,
	}
	attachmentPlexStop := slack.AttachmentField{
		Value: `✅ Stop Plex:\n` + "`/plex stop`",
		Short: false,
	}
	attachmentPlexStatus := slack.AttachmentField{
		Value: `✅ Get Plex status:\n` + "`/plex status`",
		Short: false,
	}
	attachmentPlexToken := slack.AttachmentField{
		Value: `✅ Get Plex token:\n` + "`/plex token`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color:      "#C40203",
		Title:      "👋 Need some help with" + "`/plex`?",
		TitleLink:  util.TorrentDownloadHookURL,
		Fields: []slack.AttachmentField{
			attachmentPlexDisableMonitoring,
			attachmentPlexEnableMonitoring,
			attachmentPlexStart,
			attachmentPlexStop,
			attachmentPlexStatus,
			attachmentPlexToken,
		},
		Footer:     "Torrent Client",
		FooterIcon: "https://emoji.slack-edge.com/TD00VE755/transmission/51fa8bddc5425861.png",
		Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
	}
	attachments = append(attachments, attachment)
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending download links: [%s]\n", err)
	}
}

func SendShowsHelp(slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, 2)
	attachmentRequestShows := slack.AttachmentField{
		Value: `✅ Request shows via ombi:\n` + "`/shows-request`",
		Short: false,
	}
	attachmentSearchShows := slack.AttachmentField{
		Value: `✅ Search shows using sonarr client:\n` + "`/shows-search`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color:     "#C40203",
		Title:     "👋 Need some help with" + "`/shows-request` or `/shows-search`",
		TitleLink: util.TorrentDownloadHookURL,
		Fields: []slack.AttachmentField{
			attachmentRequestShows,
			attachmentSearchShows,
		},
	}
	attachments = append(attachments, attachment)
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending download links: [%s]\n", err)
	}
}

func SendTorrentHelp(slackIncomingHookURL string) {
	var attachments = make([]slack.Attachment, 0, 4)
	attachmentTorrentSearch := slack.AttachmentField{
		Value: `✅ Search for a movie using the given keyword(s):\n` + "`/torrent-search`",
		Short: false,
	}
	attachmentTorrentStatus := slack.AttachmentField{
		Value: `✅ Show the status the shows and movies being downloaded: \n` + "`/torrent-status`",
		Short: false,
	}
	attachmentTorrentMoviesDownload := slack.AttachmentField{
		Value: `✅ Download a movie using Transmission client:\n` + "`/torrent-movies-download`",
		Short: false,
	}
	attachmentTorrentShowsDownload := slack.AttachmentField{
		Value: `✅ Download a show using Transmission client:\n` + "`/torrent-shows-download`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color:     "#C40203",
		Title:     "👋 Need some help with" + "`/torrent-search`, `/torrent-status`, `/torrent-movies-download` or ``/torrent-shows-download`?",
		TitleLink: util.TorrentDownloadHookURL,
		Fields: []slack.AttachmentField{
			attachmentTorrentSearch,
			attachmentTorrentStatus,
			attachmentTorrentMoviesDownload,
			attachmentTorrentShowsDownload,
		},
	}
	attachments = append(attachments, attachment)
	message := &slack.WebhookMessage{
		Attachments: attachments,
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending download links: [%s]\n", err)
	}
}
