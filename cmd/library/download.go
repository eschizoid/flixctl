package library

import (
	"io/ioutil"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	"github.com/eschizoid/flixctl/models"
	"github.com/spf13/cobra"
)

var DownloadLibraryCmd = &cobra.Command{
	Use:   "download",
	Short: "To Download Movies Or Shows",
	Long:  "to download movies or shows from the library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		jobOutputOutput := glacierService.GetJobOutput(svc, jobID)
		defer jobOutputOutput.Body.Close()
		var response, err = ioutil.ReadAll(jobOutputOutput.Body)
		ShowError(err)
		retrievedFile, err := ioutil.TempFile("/plex/glacier", "movie.*.zip")
		ShowError(err)
		glacierService.Unzip(retrievedFile.Name())
		err = ioutil.WriteFile(retrievedFile.Name(), response, 0644)
		glacierService.Cleanup([]string{retrievedFile.Name()})
		ShowError(err)
		close(shutdownCh)
	},
}

type InventoryRetrieve struct {
	InventoryDate string                    `json:"InventoryDate"`
	VaultARN      string                    `json:"VaultARN"`
	ArchiveList   []models.InventoryArchive `json:"ArchiveList"`
}
