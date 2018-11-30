package library

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const ec2StatusRunning = "Running"

var (
	RootLibraryCmd = &cobra.Command{
		Use:   "library",
		Short: "To Control Media Library",
	}
	fileName             string
	jobID                string
	retrievalType        string
	query                string
	slackIncomingHookURL string
	slackNotification    string
)

var (
	_ = func() struct{} {
		ArchiveLibraryCmd.Flags().StringVarP(&fileName,
			"file",
			"f",
			"",
			"the location of the movie or show to archive",
		)
		SyncLibraryCmd.Flags().StringVarP(&query,
			"query",
			"q",
			"",
			"name of the movie or show to try to find",
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
		JobsLibraryCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_LIBRARY_INCOMING_HOOK_URL"),
			"slack channel to notify of the plex event",
		)
		JobsLibraryCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
		)
		JobsLibraryCmd.Flags().StringVarP(&retrievalType,
			"type",
			"t",
			"",
			"to retrieve archived catalogue or a list of archives(movie, show)",
		)
		RootLibraryCmd.AddCommand(ArchiveLibraryCmd, InitiateLibraryCmd, RetrieveLibraryCmd, JobsLibraryCmd, SyncLibraryCmd, CatalogueLibraryCmd)
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

func ShowError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
