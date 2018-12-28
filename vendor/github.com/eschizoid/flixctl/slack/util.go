package slack

import (
	"fmt"
	"os"
	"time"
)

var (
	TorrentDownloadHookURL = fmt.Sprintf("https://%s:9000/hooks/%s", os.Getenv("FLIXCTL_HOST"), "torrent-download")
	RetrieveJobHookURL     = fmt.Sprintf("https://%s:9000/hooks/%s", os.Getenv("FLIXCTL_HOST"), "retrieve-job")
)

func GetTimeStamp() int64 {
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println(err)
	}
	return time.Now().In(location).Unix()
}
