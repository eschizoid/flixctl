package slack

import (
	"fmt"
	"os"
	"time"
)

var (
	TorrentDownloadHookURL = fmt.Sprintf("%s/%s", os.Getenv("HOOKS_URL"), "torrent-download")
	RetrieveJobHookURL     = fmt.Sprintf("%s/%s", os.Getenv("HOOKS_URL"), "retrieve-job")
)

func GetTimeStamp() int64 {
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println(err)
	}
	return time.Now().In(location).Unix()
}
