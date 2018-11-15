package library

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var RootLibraryCmd = &cobra.Command{
	Use:   "library",
	Short: "To Control Media Library",
}

var retrievalType string
var fileName string
var archiveID string
var jobID string

var (
	_ = func() struct{} {
		ArchiveLibraryCmd.Flags().StringVarP(&fileName,
			"file",
			"f",
			"",
			"the location of the movie or show to archive",
		)
		InitiateLibraryCmd.Flags().StringVarP(&retrievalType,
			"type",
			"t",
			"",
			"to retrieve archived library catalogue or a list of archives(movie, show)",
		)
		InitiateLibraryCmd.Flags().StringVarP(&archiveID,
			"archive-id",
			"i",
			"",
			"to archive id to retrieve",
		)
		RetrieveLibraryCmd.Flags().StringVarP(&jobID,
			"job-id",
			"i",
			"",
			"the job id to start for retrieving a movie or a show",
		)
		RetrieveLibraryCmd.Flags().StringVarP(&fileName,
			"file",
			"f",
			"",
			"the location to retrieve the movie of the show",
		)
		RootLibraryCmd.AddCommand(ArchiveLibraryCmd, InitiateLibraryCmd, RetrieveLibraryCmd, JobsLibraryCmd)
		return struct{}{}
	}()
)

func Indicator(shutdownCh <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case <-shutdownCh:
			return
		}
	}
}
