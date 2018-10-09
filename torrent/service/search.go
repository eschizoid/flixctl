package service

import (
	"fmt"
	"regexp"
	"time"

	"github.com/juliensalinas/torrengo/otts"
	"github.com/juliensalinas/torrengo/td"
	"github.com/juliensalinas/torrengo/tpb"
)

const (
	TorrentDownloadsKey = "td"
	ThePirateBayKey     = "tpb"
	OttsKey             = "otts"
)

type TorrentResult struct {
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

type TorrentSearch struct {
	In              string
	Out             []TorrentResult
	SourcesToLookup []string
}

var Timeout = time.Duration(15000 * 1000 * 1000)

var TdTorListCh = make(chan []TorrentResult)
var TpbTorListCh = make(chan []TorrentResult)
var OttsTorListCh = make(chan []TorrentResult)

var TdSearchErrCh = make(chan error)
var TpbSearchErrCh = make(chan error)
var OttsSearchErrCh = make(chan error)

var sources = map[string]string{
	TorrentDownloadsKey: "Torrent Downloads",
	ThePirateBayKey:     "The Pirate Bay",
	OttsKey:             "1337x",
}

var regex = regexp.MustCompile("[[:^ascii:]]")

func Search(search *TorrentSearch) { //nolint:gocyclo

	for _, source := range search.SourcesToLookup {
		switch source {
		case TorrentDownloadsKey:
			go func() {
				tdTorrents, err := td.Lookup(search.In, Timeout)
				if err != nil {
					TdSearchErrCh <- err
					return
				}
				var torrentList []TorrentResult
				for _, tdTorrent := range tdTorrents {
					_, magnet, _ := td.ExtractTorAndMag(tdTorrent.DescURL, Timeout)
					result := TorrentResult{
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
				var torrentList []TorrentResult
				for _, tpbTorrent := range tpbTorrents {
					result := TorrentResult{
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
				var torrentList []TorrentResult
				for _, ottsTorrent := range ottsTorrents {
					magnet, _ := otts.ExtractMag(ottsTorrent.DescURL, Timeout)
					result := TorrentResult{
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

func Merge(search *TorrentSearch) [3]error { //nolint:gocyclo
	var tdSearchErr, tpbSearchErr, ottsSearchErr error

	// Gather all goroutines results
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
