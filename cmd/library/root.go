package library

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const ec2StatusRunning = "Running"
const ec2StatusStopped = "Stopped"

var (
	RootLibraryCmd = &cobra.Command{
		Use:   "library",
		Short: "To Control Media Library",
	}
	archiveFilter           string
	archiveID               string
	enableBatchUpload       string
	jobID                   string
	jobFilter               string
	maxUploadItems          string
	slackIncomingHookURL    string
	slackNotification       string
	enableLibrarySync       string
	sourceFile              string
	targetFile              string
	awsResourceTagNameValue = os.Getenv("AWS_RESOURCE_TAG_NAME_VALUE")
)

var (
	_ = func() struct{} {
		CatalogueLibraryCmd.Flags().StringVarP(&archiveFilter,
			"filter",
			"f",
			"",
			"the optional filter to apply when retrieving the catalogue",
		)
		DeleteArchiveLibraryCmd.Flags().StringVarP(&archiveID,
			"archive-id",
			"i",
			"",
			"the id of the archive to be deleted",
		)
		DownloadLibraryCmd.Flags().StringVarP(&jobID,
			"job-id",
			"i",
			"",
			"the optional job id for retrieving glacier archive",
		)
		DownloadLibraryCmd.Flags().StringVarP(&targetFile,
			"target-file",
			"f",
			"",
			"where to download the file",
		)
		InitiateLibraryCmd.Flags().StringVarP(&archiveID,
			"archive-id",
			"i",
			"",
			"if provided, will attempt an archive retrieval instead of an inventory retrieval",
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
			os.Getenv("SLACK_REQUESTS_HOOK_URL"),
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
			os.Getenv("SLACK_REQUESTS_HOOK_URL"),
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
		UploadLibraryCmd.Flags().StringVarP(&enableBatchUpload,
			"enable-batch-mode",
			"b",
			"",
			"optional argument to upload all the watched movies and shows.",
		)
		UploadLibraryCmd.Flags().StringVarP(&maxUploadItems,
			"max-upload-items",
			"m",
			"",
			"optional max number of items to upload when bach mode was enabled.",
		)
		UploadLibraryCmd.Flags().StringVarP(&sourceFile,
			"source-file",
			"f",
			"",
			"the source file to upload to the library.",
		)
		RootLibraryCmd.AddCommand(
			DeleteArchiveLibraryCmd,
			DownloadLibraryCmd,
			CatalogueLibraryCmd,
			InitiateLibraryCmd,
			InventoryLibraryCmd,
			JobsLibraryCmd,
			SyncLibraryCmd,
			UploadLibraryCmd,
		)

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
