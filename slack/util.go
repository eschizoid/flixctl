package slack

import (
	"fmt"
	"os"
	"time"
)

var (
	TorrentDownloadHookURL  = fmt.Sprintf("https://%s:9000/hooks/%s", os.Getenv("FLIXCTL_HOST"), "torrent-download")
	LibraryInitiateHookURL  = fmt.Sprintf("https://%s:9000/hooks/%s", os.Getenv("FLIXCTL_HOST"), "library-initiate")
	LibraryInventoryHookURL = fmt.Sprintf("https://%s:9000/hooks/%s", os.Getenv("FLIXCTL_HOST"), "library-inventory")
	LibraryDownloadHookURL  = fmt.Sprintf("https://%s:9000/hooks/%s", os.Getenv("FLIXCTL_HOST"), "library-download")
	LibraryDeleteHookURL    = fmt.Sprintf("https://%s:9000/hooks/%s", os.Getenv("FLIXCTL_HOST"), "library-delete")
)

func GetTimeStamp() int64 {
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println(err)
	}
	return time.Now().In(location).Unix()
}
