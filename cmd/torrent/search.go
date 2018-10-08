package torrent

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Jeffail/gabs"
	torrentService "github.com/eschizoid/flixctl/torrent/service"
	"github.com/spf13/cobra"
)

var SearchTorrentCmd = &cobra.Command{
	Use:   "search",
	Short: "To Search for Torrents",
	Long:  `to search for torrents using the given keywors(s)`,
	Run: func(cmd *cobra.Command, args []string) {
		torrentSearch := torrentService.TorrentSearch{
			In: keywords,
			SourcesToLookup: []string{
				torrentService.ThePirateBayKey,
				torrentService.OttsKey,
				torrentService.TorrentDownloadsKey,
			},
		}
		var err = cleanInput(torrentSearch)
		if err != nil {
			fmt.Println("Could not process your input")
		}
		torrentService.Search(&torrentSearch)
		errors := torrentService.Merge(&torrentSearch)
		if errors[0] != nil && errors[1] != nil && errors[2] != nil {
			fmt.Println("All searches returned an error.")
		}
		if len(torrentSearch.Out) == 0 {
			fmt.Println("No result found...")
		}
		choose(&torrentSearch)
		sortOut(&torrentSearch)
		resultJSON, _ := gabs.Consume(torrentSearch)
		fmt.Println(string(resultJSON.EncodeJSON()))
	},
}

func cleanInput(search torrentService.TorrentSearch) error {
	var in = strings.TrimSpace(search.In)
	if in == "" {
		return fmt.Errorf("user search should not be empty")
	}
	return nil
}

func choose(search *torrentService.TorrentSearch) {
	var test = func(s string) bool { return strings.Contains(s, quality) }
	var filteredResult []torrentService.TorrentResult
	for _, torrentResult := range search.Out {
		if test(torrentResult.Name) {
			torrentResult.Quality = quality + "p"
			filteredResult = append(filteredResult, torrentResult)
		}
	}
	search.Out = filteredResult
}

func sortOut(search *torrentService.TorrentSearch) {
	sort.Slice(search.Out, func(i, j int) bool {
		return search.Out[i].Seeders > search.Out[j].Seeders
	})
}
