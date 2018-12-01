package torrent

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	slackService "github.com/eschizoid/flixctl/slack/torrent"
	torrentService "github.com/eschizoid/flixctl/torrent"
	"github.com/spf13/cobra"
)

var SearchTorrentCmd = &cobra.Command{
	Use:   "search",
	Short: "To Search For Torrents",
	Long:  "to search for torrents using the given keyword(s)",
	Run: func(cmd *cobra.Command, args []string) {
		torrentSearch := torrentService.Search{
			In: keywords,
			SourcesToLookup: []string{
				torrentService.ThePirateBayKey,
				torrentService.OttsKey,
				torrentService.TorrentDownloadsKey,
			},
		}
		var err = cleanInput(&torrentSearch)
		if err != nil {
			fmt.Println("Could not process your input")
		}
		torrentService.SearchTorrents(&torrentSearch)
		errors := torrentService.Merge(&torrentSearch)
		if errors[0] != nil && errors[1] != nil && errors[2] != nil {
			fmt.Println("All searches returned an error")
		}
		if len(torrentSearch.Out) == 0 {
			fmt.Println("No torrents found")
		}
		choose(&torrentSearch)
		sortOut(&torrentSearch)
		out, _ := json.Marshal(torrentSearch.Out)
		fmt.Println(string(out))
		if notify, _ := strconv.ParseBool(slackNotification); notify {
			slackService.SendDownloadLinks(&torrentSearch, slackIncomingHookURL, downloadDir, true)
		}
	},
}

func cleanInput(search *torrentService.Search) error {
	var in = strings.TrimSpace(search.In)
	if in == "" {
		return fmt.Errorf("user search should not be empty")
	}
	return nil
}

func choose(search *torrentService.Search) {
	var test = func(s string) bool { return strings.Contains(s, quality) }
	var filteredResult []torrentService.Result
	for _, torrentResult := range search.Out {
		if test(torrentResult.Name) {
			torrentResult.Quality = quality + "p"
			filteredResult = append(filteredResult, torrentResult)
		}
	}
	if len(filteredResult) > 0 {
		search.Out = filteredResult
	}
}

func sortOut(search *torrentService.Search) {
	sort.Slice(search.Out, func(i, j int) bool {
		return search.Out[i].Seeders > search.Out[j].Seeders
	})
}
