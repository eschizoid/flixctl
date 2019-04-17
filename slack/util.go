package slack

import (
	"fmt"
	"os"
	"time"
)

var (
	TorrentDownloadHookURL        = fmt.Sprintf("https://%s/prod/torrent/download?lambda=torrent-executor", os.Getenv("FLIXCTL_HOST"))
	LibraryInitiateArchiveHookURL = fmt.Sprintf("https://%s/prod/torrent/download?lambda=library-executor", os.Getenv("FLIXCTL_HOST"))
	LibraryInventoryHookURL       = fmt.Sprintf("https://%s/prod/torrent/download?lambda=library-executor", os.Getenv("FLIXCTL_HOST"))
	LibraryDownloadHookURL        = fmt.Sprintf("https://%s/prod/torrent/download?lambda=library-executor", os.Getenv("FLIXCTL_HOST"))
	LibraryDeleteHookURL          = fmt.Sprintf("https://%s/prod/torrent/download?lambda=library-executor", os.Getenv("FLIXCTL_HOST"))
	SigningSecret                 = os.Getenv("SLACK_SIGNING_SECRET")
)

func GetTimeStamp() int64 {
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println(err)
	}
	return time.Now().In(location).Unix()
}
