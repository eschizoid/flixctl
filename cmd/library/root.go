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
	jobID                string
	jobFilter            string
	slackIncomingHookURL string
	slackNotification    string
	enableLibrarySync    string
	sourceFile           string
)

var (
	_ = func() struct{} {
		DownloadLibraryCmd.Flags().StringVarP(&jobID,
			"job-id",
			"i",
			"",
			"the job id for retrieving glacier archive inventory",
		)
		InventoryLibraryCmd.Flags().StringVarP(&enableLibrarySync,
			"enable-sync",
			"e",
			"",
			"optional argument to sync glacier archive inventory with internal library",
		)
		InventoryLibraryCmd.Flags().StringVarP(&jobID,
			"job-id",
			"i",
			"",
			"the optional job id for retrieving glacier archive inventory",
		)
		InventoryLibraryCmd.Flags().StringVarP(&slackIncomingHookURL,
			"slack-notification-channel",
			"s",
			os.Getenv("SLACK_LIBRARY_INCOMING_HOOK_URL"),
			"slack channel to notify of the plex event",
		)
		InventoryLibraryCmd.Flags().StringVarP(&slackNotification,
			"slack-notification",
			"n",
			os.Getenv("SLACK_NOTIFICATION"),
			"if true, will try to notify to a slack channel",
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
		JobsLibraryCmd.Flags().StringVarP(&jobFilter,
			"filter",
			"f",
			"",
			"to filter the list of jobs.",
		)
		UploadLibraryCmd.Flags().StringVarP(&sourceFile,
			"source-file",
			"f",
			"",
			"the source file to upload to the library.",
		)
		RootLibraryCmd.AddCommand(UploadLibraryCmd, InitiateLibraryCmd, DownloadLibraryCmd, InventoryLibraryCmd,
			JobsLibraryCmd, SyncLibraryCmd, CatalogueLibraryCmd)
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
