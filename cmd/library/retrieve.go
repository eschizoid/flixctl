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
	Short: "To Retrieve A Movie Or Show",
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
		var file *os.File
		if retrievalType == "InventoryRetrieval" {
			file, err = ioutil.TempFile(os.TempDir(), "inventory.*.json")
			ShowError(err)
		} else if retrievalType == "ArchiveRetrieval" {
			file, err = ioutil.TempFile(os.TempDir(), "movie.*.zip")
			ShowError(err)
		}
		err = ioutil.WriteFile(file.Name(), part, 0644)
		ShowError(err)
		jsonString, _ := json.Marshal(getJobOutputOutput)
		glacierService.Cleanup([]string{file.Name()})
		fmt.Println("\n" + string(jsonString))
		close(shutdownCh)
	},
}

func ShowError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
