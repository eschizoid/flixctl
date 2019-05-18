package lambda

import (
	"encoding/json"
	"fmt"
	"strconv"

	util "github.com/eschizoid/flixctl/slack"
	"github.com/nlopes/slack"
)

const TitleLink = "https://github.com/eschizoid/flixctl/blob/master/README.adoc"

var attachmentLearnMore = slack.Attachment{
	Color:     "#f2f2f2",
	TitleLink: TitleLink,
	Title:     "Need more help?",
}

func SendAdminHelp(slackIncomingHookURL string) {
	attachmentRenewCerts := slack.AttachmentField{
		Value: "✅ Renew TLS certificates for all plex related services:\n`/admin renew-certs`",
		Short: false,
	}
	attachmentRestartServices := slack.AttachmentField{
		Value: "✅ Restart all Plex related services:\n`/admin restart-services`",
		Short: false,
	}
	attachmentSlackToken := slack.AttachmentField{
		Value: "✅ Get a Slack oauth token for a given client id:\n`/admin slack-token`",
		Short: false,
	}
	attachmentSlackPurge := slack.AttachmentField{
		Value: "✅ Purge Slack messages from all channels:\n`/admin slack-purge`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color: "#f2f2f2",
		Text:  "👋 Need some help with `/admin`?",
		Fields: []slack.AttachmentField{
			attachmentRenewCerts,
			attachmentRestartServices,
			attachmentSlackToken,
			attachmentSlackPurge,
		},
		MarkdownIn: []string{"text", "fields"},
	}
	message := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			attachment,
			attachmentLearnMore,
		},
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending admin help: [%s]\n", err)
	}
}

func SendMoviesHelp(slackIncomingHookURL string) {
	attachmentRequestMovie := slack.AttachmentField{
		Value: "✅ Request movies via Ombi:\n`/movies-request`",
		Short: false,
	}
	attachmentSearchMovie := slack.AttachmentField{
		Value: "✅ Search movies using Radarr client:\n`/movies-search`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color: "#f2f2f2",
		Text:  "👋 Need some help with `/movies-request` or `/movies-search`?",
		Fields: []slack.AttachmentField{
			attachmentRequestMovie,
			attachmentSearchMovie,
		},
		MarkdownIn: []string{"text", "fields"},
		Ts:         json.Number(strconv.FormatInt(util.GetTimeStamp(), 10)),
	}
	message := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			attachment,
			attachmentLearnMore,
		},
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending movies help: [%s]\n", err)
	}
}

func SendPlexHelp(slackIncomingHookURL string) {
	attachmentPlexDisableMonitoring := slack.AttachmentField{
		Value: "✅ Enable Plex monitoring:\n`/plex disable-monitoring`",
		Short: false,
	}
	attachmentPlexEnableMonitoring := slack.AttachmentField{
		Value: "✅ Disable Plex monitoring:\n`/plex enable-monitoring`",
		Short: false,
	}
	attachmentPlexStart := slack.AttachmentField{
		Value: "✅ Start Plex:\n`/plex start`",
		Short: false,
	}
	attachmentPlexStop := slack.AttachmentField{
		Value: "✅ Stop Plex:\n`/plex stop`",
		Short: false,
	}
	attachmentPlexStatus := slack.AttachmentField{
		Value: "✅ Get Plex status:\n`/plex status`",
		Short: false,
	}
	attachmentPlexToken := slack.AttachmentField{
		Value: "✅ Get Plex token:\n`/plex token`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color: "f2f2f2 ",
		Text:  "👋 Need some help with `/plex`?",
		Fields: []slack.AttachmentField{
			attachmentPlexDisableMonitoring,
			attachmentPlexEnableMonitoring,
			attachmentPlexStart,
			attachmentPlexStop,
			attachmentPlexStatus,
			attachmentPlexToken,
		},
		MarkdownIn: []string{"text", "fields"},
	}
	message := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			attachment,
			attachmentLearnMore,
		},
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending plex help: [%s]\n", err)
	}
}

func SendShowsHelp(slackIncomingHookURL string) {
	attachmentRequestShows := slack.AttachmentField{
		Value: "✅ Request shows via Ombi:\n`/shows-request`",
		Short: false,
	}
	attachmentSearchShows := slack.AttachmentField{
		Value: "✅ Search shows using Sonarr client:\n`/shows-search`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color: "#f2f2f2",
		Text:  "👋 Need some help with `/shows-request` or `/shows-search`?",
		Fields: []slack.AttachmentField{
			attachmentRequestShows,
			attachmentSearchShows,
		},
		MarkdownIn: []string{"text", "fields"},
	}
	message := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			attachment,
			attachmentLearnMore,
		},
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending shows help: [%s]\n", err)
	}
}

func SendTorrentHelp(slackIncomingHookURL string) {
	attachmentTorrentSearch := slack.AttachmentField{
		Value: "✅ Search for a movie using the given keyword(s):\n`/torrent-search`",
		Short: false,
	}
	attachmentTorrentStatus := slack.AttachmentField{
		Value: "✅ Show the status the shows and movies being downloaded:\n`/torrent-status`",
		Short: false,
	}
	attachmentTorrentMoviesDownload := slack.AttachmentField{
		Value: "✅ Download a movie using Transmission client:\n`/torrent-movies-download`",
		Short: false,
	}
	attachmentTorrentShowsDownload := slack.AttachmentField{
		Value: "✅ Download a show using Transmission client:\n`/torrent-shows-download`",
		Short: false,
	}
	attachment := slack.Attachment{
		Color: "#f2f2f2",
		Text:  "👋 Need some help with `/torrent-search`, `/torrent-status`, `/torrent-movies-download` or `/torrent-shows-download`?",
		Fields: []slack.AttachmentField{
			attachmentTorrentSearch,
			attachmentTorrentStatus,
			attachmentTorrentMoviesDownload,
			attachmentTorrentShowsDownload,
		},
		MarkdownIn: []string{"text", "fields"},
	}
	message := &slack.WebhookMessage{
		Attachments: []slack.Attachment{
			attachment,
			attachmentLearnMore,
		},
	}
	err := slack.PostWebhook(slackIncomingHookURL, message)
	if err != nil {
		fmt.Printf("Error while sending torrent help: [%s]\n", err)
	}
}
