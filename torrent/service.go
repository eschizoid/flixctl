package torrent

import (
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	"github.com/juliensalinas/torrengo/otts"
	"github.com/juliensalinas/torrengo/td"
	"github.com/juliensalinas/torrengo/tpb"
)

const (
	ec2RunningStatus     = "Running"
	transmissionHostPort = "marianoflix.duckdns.org:9091"
	TorrentDownloadsKey  = "td"
	ThePirateBayKey      = "tpb"
	OttsKey              = "otts"
)

type Result struct {
	FileURL    string
	Magnet     string
	DescURL    string
	Name       string
	Size       string
	Quality    string
	Seeders    int
	Leechers   int
	UploadDate string
	Source     string
	FilePath   string
}

type Search struct {
	In              string
	Out             []Result
	SourcesToLookup []string
}

var Timeout = time.Duration(15000 * 1000 * 1000)

var TdTorListCh = make(chan []Result)
var TpbTorListCh = make(chan []Result)
var OttsTorListCh = make(chan []Result)

var TdSearchErrCh = make(chan error)
var TpbSearchErrCh = make(chan error)
var OttsSearchErrCh = make(chan error)

var Sources = map[string]string{
	TorrentDownloadsKey: "Torrent Downloads",
	ThePirateBayKey:     "The Pirate Bay",
	OttsKey:             "1337x",
}

var AwsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
	SharedConfigState: sess.SharedConfigEnable,
}))

var SVC = ec2.New(AwsSession, AwsSession.Config)

var InstanceID = ec2Service.FetchInstanceID(SVC, "plex")

var Regex = regexp.MustCompile("[[:^ascii:]]")

func SearchTorrents(search *Search) { //nolint:gocyclo

	for _, source := range search.SourcesToLookup {
		switch source {
		case TorrentDownloadsKey:
			go func() {
				tdTorrents, err := td.Lookup(search.In, Timeout)
				if err != nil {
					TdSearchErrCh <- err
					return
				}
				var torrentList []Result
				for _, tdTorrent := range tdTorrents {
					_, magnet, _ := td.ExtractTorAndMag(tdTorrent.DescURL, Timeout)
					result := Result{
						DescURL:  tdTorrent.DescURL,
						Magnet:   magnet,
						Name:     tdTorrent.Name,
						Size:     tdTorrent.Size,
						Leechers: tdTorrent.Leechers,
						Seeders:  tdTorrent.Seeders,
						Source:   TorrentDownloadsKey,
					}
					torrentList = append(torrentList, result)
				}
				TdTorListCh <- torrentList
			}()
		case ThePirateBayKey:
			go func() {
				tpbTorrents, err := tpb.Lookup(search.In, Timeout)
				if err != nil {
					TpbSearchErrCh <- err
					return
				}
				var torrentList []Result
				for _, tpbTorrent := range tpbTorrents {
					result := Result{
						Magnet:     tpbTorrent.Magnet,
						Name:       Regex.ReplaceAllLiteralString(tpbTorrent.Name, " "),
						Size:       Regex.ReplaceAllLiteralString(tpbTorrent.Size, " "),
						UploadDate: Regex.ReplaceAllLiteralString(tpbTorrent.UplDate, " "),
						Leechers:   tpbTorrent.Leechers,
						Seeders:    tpbTorrent.Seeders,
						Source:     ThePirateBayKey,
					}
					torrentList = append(torrentList, result)
				}
				TpbTorListCh <- torrentList
			}()
		case OttsKey:
			go func() {
				ottsTorrents, err := otts.Lookup(search.In, Timeout)
				if err != nil {
					OttsSearchErrCh <- err
					return
				}
				var torrentList []Result
				for _, ottsTorrent := range ottsTorrents {
					magnet, _ := otts.ExtractMag(ottsTorrent.DescURL, Timeout)
					result := Result{
						DescURL:    ottsTorrent.DescURL,
						Magnet:     magnet,
						Name:       Regex.ReplaceAllLiteralString(ottsTorrent.Name, " "),
						Size:       Regex.ReplaceAllLiteralString(ottsTorrent.Size, " "),
						UploadDate: Regex.ReplaceAllLiteralString(ottsTorrent.UplDate, " "),
						Leechers:   ottsTorrent.Leechers,
						Seeders:    ottsTorrent.Seeders,
						Source:     OttsKey,
					}
					torrentList = append(torrentList, result)
				}
				OttsTorListCh <- torrentList
			}()
		}
	}
}

func Merge(search *Search) [3]error { //nolint:gocyclo
	var tdSearchErr, tpbSearchErr, ottsSearchErr error

	for _, source := range search.SourcesToLookup {
		switch source {
		case TorrentDownloadsKey:
			select {
			case tdSearchErr = <-TdSearchErrCh:
				fmt.Printf("An error occured during search on %v\n", Sources["td"])
			case tdTorList := <-TdTorListCh:
				search.Out = append(search.Out, tdTorList...)
			}
		case ThePirateBayKey:
			select {
			case tpbSearchErr = <-TpbSearchErrCh:
				fmt.Printf("An error occured during search on %v\n", Sources["tpb"])
			case tpbTorList := <-TpbTorListCh:
				search.Out = append(search.Out, tpbTorList...)
			}
		case OttsKey:
			select {
			case ottsSearchErr = <-OttsSearchErrCh:
				fmt.Printf("An error occured during search on %v\n", Sources["otts"])
			case ottsTorList := <-OttsTorListCh:
				search.Out = append(search.Out, ottsTorList...)
			}
		}
	}
	errors := [3]error{tdSearchErr, tpbSearchErr, ottsSearchErr}
	return errors
}

func Status() string {
	var torrentStatus string
	ec2status := ec2Service.Status(SVC, InstanceID)
	if ec2status == ec2RunningStatus {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "transmission-remote",
			transmissionHostPort,
			"--authenv",
			"--torrent=active",
			"--list")
		out, err := cmd.CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			torrentStatus = "Command timed out"
			fmt.Printf("Could not list torrents being downloaded: [%s]\n", err)
		} else {
			torrentStatus = string(out)
		}
	} else {
		torrentStatus = "Plex Stopped"
	}
	return torrentStatus
}

func TriggerDownload(envMagnetLink string, argMagnetLink string, envDownloadDir string) {
	if envMagnetLink == "" {
		// coming from flixctl
		downloadTorrent(argMagnetLink, envDownloadDir)
	} else {
		// coming from web-hook
		decodedEnvMagnetLink, err := base64.StdEncoding.DecodeString(envMagnetLink)
		if err != nil {
			fmt.Printf("Could not decode the magnet link: [%s]\n", err)
		}
		downloadTorrent(string(decodedEnvMagnetLink), envDownloadDir)
	}
}

func downloadTorrent(magnet string, downloadDir string) {
	status := ec2Service.Status(SVC, InstanceID)
	if status == ec2RunningStatus {
		transmission := exec.Command("transmission-remote",
			transmissionHostPort,
			"--authenv",
			"--add",
			fmt.Sprintf("--download-dir=%s", downloadDir),
			magnet)
		err := transmission.Start()
		if err != nil {
			fmt.Println("Could not download torrent using the given magnet link")
		}
	}
}
