package library

import (
	"os"

	"github.com/eschizoid/flixctl/library"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/spf13/cobra"
)

var SyncLibraryCmd = &cobra.Command{
	Use:   "sync",
	Short: "To Sync Plex DB",
	Long:  "to sync plex db with internal flixctl db.",
	Run: func(cmd *cobra.Command, args []string) {
		token := os.Getenv("PLEX_TOKEN")
		movies := library.GetPlexMovies(token)
		for _, movie := range movies.MediaContainer.Metadata {
			err := libraryService.SaveMovie(movie)
			ShowError(err)
		}
	},
}
