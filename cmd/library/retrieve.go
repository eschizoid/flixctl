package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	"github.com/spf13/cobra"
)

var RetrieveLibraryCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "To Retrieve Movie Or Show",
	Long:  "to retrieve a movie or show from the media library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		getJobOutputOutput := glacierService.GetJobOutput(svc, jobID)
		defer getJobOutputOutput.Body.Close()
		part, err := ioutil.ReadAll(getJobOutputOutput.Body)
		ShowError(err)
		var zipFileName = writeFile(part)
		glacierService.Unzip(zipFileName)
		glacierService.Cleanup([]string{zipFileName})
		jsonString, _ := json.Marshal(getJobOutputOutput)
		fmt.Println("\n" + string(jsonString))
		close(shutdownCh)
	},
}

func writeFile(part []byte) string {
	var err error
	var zipFile *os.File
	if retrievalType == "InventoryRetrieval" {
		zipFile, err = ioutil.TempFile(os.TempDir(), "inventory.*.json")
	} else if retrievalType == "ArchiveRetrieval" {
		zipFile, err = ioutil.TempFile(os.TempDir(), "movie.*.zip")
	}
	ShowError(err)
	zipFileName := zipFile.Name()
	err = ioutil.WriteFile(zipFileName, part, 0644)
	ShowError(err)
	return zipFileName
}
