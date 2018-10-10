package torrent

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	"github.com/eschizoid/flixctl/cmd/plex"
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

var sources = map[string]string{
	TorrentDownloadsKey: "Torrent Downloads",
	ThePirateBayKey:     "The Pirate Bay",
	OttsKey:             "1337x",
}

var regex = regexp.MustCompile("[[:^ascii:]]")

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
						Name:       regex.ReplaceAllLiteralString(tpbTorrent.Name, " "),
						Size:       regex.ReplaceAllLiteralString(tpbTorrent.Size, " "),
						UploadDate: regex.ReplaceAllLiteralString(tpbTorrent.UplDate, " "),
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
						Name:       regex.ReplaceAllLiteralString(ottsTorrent.Name, " "),
						Size:       regex.ReplaceAllLiteralString(ottsTorrent.Size, " "),
						UploadDate: regex.ReplaceAllLiteralString(ottsTorrent.UplDate, " "),
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
				fmt.Printf("An error occured during search on %v\n", sources["td"])
			case tdTorList := <-TdTorListCh:
				search.Out = append(search.Out, tdTorList...)
			}
		case ThePirateBayKey:
			select {
			case tpbSearchErr = <-TpbSearchErrCh:
				fmt.Printf("An error occured during search on %v\n", sources["tpb"])
			case tpbTorList := <-TpbTorListCh:
				search.Out = append(search.Out, tpbTorList...)
			}
		case OttsKey:
			select {
			case ottsSearchErr = <-OttsSearchErrCh:
				fmt.Printf("An error occured during search on %v\n", sources["otts"])
			case ottsTorList := <-OttsTorListCh:
				search.Out = append(search.Out, ottsTorList...)
			}
		}
	}
	errors := [...]error{tdSearchErr, tpbSearchErr, ottsSearchErr}
	return errors
}

func Status() string {
	var torrentStatus string
	ec2status := ec2Service.Status(plex.Session, plex.InstanceID)
	if ec2status == ec2RunningStatus {
		out, err := exec.Command("transmission-remote",
			transmissionHostPort,
			"--authenv",
			"--torrent=all",
			"--list").CombinedOutput()
		if err != nil {
			fmt.Printf("Could not list torrents being downloaded: [%s]\n", err)
		}
		torrentStatus = string(out)
		fmt.Println(torrentStatus)
	}
	return torrentStatus
}

func TriggerDownload(envMagnetLink string, argMagnetLink string) {
	if envMagnetLink == "" {
		// coming from flixctl
		downloadTorrent(argMagnetLink)
	} else {
		// coming from web-hook
		decodedEnvMagnetLink, err := base64.StdEncoding.DecodeString(envMagnetLink)
		if err != nil {
			fmt.Printf("Could not decode the magnet link: [%s]\n", err)
		}
		downloadTorrent(string(decodedEnvMagnetLink))
	}
}

func downloadTorrent(magnet string) {
	status := ec2Service.Status(plex.Session, plex.InstanceID)
	if status == ec2RunningStatus {
		transmission := exec.Command("transmission-remote",
			transmissionHostPort,
			"--authenv",
			"--add",
			magnet)
		err := transmission.Start()
		if err != nil {
			fmt.Println("Could not download torrent using the given magnet link")
		}
	}
}
