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
	Short: "To Retrieve Movies Or Shows",
	Long:  "to retrieve movies or shows from the library.",
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
		writeFile(part)
		//glacierService.Cleanup([]string{fileName})
		jsonString, _ := json.Marshal(getJobOutputOutput)
		fmt.Println("\n" + string(jsonString))
		close(shutdownCh)
	},
}

func writeFile(part []byte) string {
	var err error
	var retrievedFile *os.File
	if retrievalType == "InventoryRetrieval" {
		retrievedFile, err = ioutil.TempFile(os.TempDir(), "inventory.*.json")
	} else if retrievalType == "ArchiveRetrieval" {
		retrievedFile, err = ioutil.TempFile(os.TempDir(), "movie.*.zip")
		glacierService.Unzip(retrievedFile.Name())
	}
	ShowError(err)
	fileName := retrievedFile.Name()
	err = ioutil.WriteFile(fileName, part, 0644)
	ShowError(err)
	return fileName
}
