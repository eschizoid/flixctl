package torrent

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hekmon/transmissionrpc"
	"github.com/juliensalinas/torrengo/otts"
	"github.com/juliensalinas/torrengo/td"
	"github.com/juliensalinas/torrengo/tpb"
)

const (
	ec2RunningStatus    = "Running"
	transmissionHost    = "marianoflix.duckdns.org"
	TorrentDownloadsKey = "td"
	ThePirateBayKey     = "tpb"
	OttsKey             = "otts"
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

var (
	Regex           = regexp.MustCompile("[[:^ascii:]]")
	transmission, _ = transmissionrpc.New(
		transmissionHost,
		strings.Split(os.Getenv("TR_AUTH"), ":")[0],
		strings.Split(os.Getenv("TR_AUTH"), ":")[1],
		&transmissionrpc.AdvancedConfig{
			HTTPS: true,
			Port:  443,
		})
)

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
				fmt.Printf("An error occurred during search on %v\n", Sources["td"])
			case tdTorList := <-TdTorListCh:
				search.Out = append(search.Out, tdTorList...)
			}
		case ThePirateBayKey:
			select {
			case tpbSearchErr = <-TpbSearchErrCh:
				fmt.Printf("An error occurred during search on %v\n", Sources["tpb"])
			case tpbTorList := <-TpbTorListCh:
				search.Out = append(search.Out, tpbTorList...)
			}
		case OttsKey:
			select {
			case ottsSearchErr = <-OttsSearchErrCh:
				fmt.Printf("An error occurred during search on %v\n", Sources["otts"])
			case ottsTorList := <-OttsTorListCh:
				search.Out = append(search.Out, ottsTorList...)
			}
		}
	}
	errors := [3]error{tdSearchErr, tpbSearchErr, ottsSearchErr}
	return errors
}

func Status(ec2status string) []transmissionrpc.Torrent {
	var torrents []transmissionrpc.Torrent
	if ec2status == ec2RunningStatus {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		response, err := transmission.TorrentGetAll()
		if err != nil {
			panic(err)
		}
		for _, torrent := range response {
			if torrent.Files[0].Name != "" && *torrent.PercentDone < 1.00 {
				torrents = append(torrents, *torrent)
			}
		}
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("Could not list torrents being downloaded")
		}
	}
	return torrents
}

func TriggerDownload(magnetLink string, downloadDir string, ec2status string) *transmissionrpc.Torrent {
	var err error
	var torrent *transmissionrpc.Torrent
	if ec2status == ec2RunningStatus {
		torrent, err = transmission.TorrentAdd(&transmissionrpc.TorrentAddPayload{
			DownloadDir: &downloadDir,
			Filename:    &magnetLink,
		})
		if err != nil {
			fmt.Printf("Could not download torrent using the given magnet link: [%s]", err)
		}
	}
	return torrent
}
