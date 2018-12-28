package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/eschizoid/flixctl/models"
	"github.com/spf13/cobra"
)

const (
	GlacierInventoryRetrieval string = "InventoryRetrieval"
	GlacierArchiveRetrieval   string = "ArchiveRetrieval"
)

var RetrieveLibraryCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "To Retrieve Movies Or Shows",
	Long:  "to retrieve movies or shows from the library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		if retrievalType == GlacierInventoryRetrieval {
			processInventoryRetrieval()
		} else if retrievalType == GlacierArchiveRetrieval {
			processArchiveRetrieval()
		} else {
			panic("Unknown glacier retrieval type")
		}
		close(shutdownCh)
	},
}

func processInventoryRetrieval() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	if jobID != "" {
		svc := glacier.New(awsSession)
		jobOutputOutput := glacierService.GetJobOutput(svc, jobID)
		defer jobOutputOutput.Body.Close()
		response, err := ioutil.ReadAll(jobOutputOutput.Body)
		ShowError(err)
		var inventoryRetrieve = new(InventoryRetrieve)
		err = json.Unmarshal(response, &inventoryRetrieve)
		ShowError(err)
		for _, archive := range inventoryRetrieve.ArchiveList {
			err = libraryService.SaveGlacierInventoryArchive(archive)
			ShowError(err)
		}
	}
	glacierArchives, err := libraryService.GetGlacierInventoryArchives()
	ShowError(err)
	jsonString, err := json.Marshal(glacierArchives)
	ShowError(err)
	fmt.Println("\n" + string(jsonString))
}

func processArchiveRetrieval() {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	svc := glacier.New(awsSession)
	jobOutputOutput := glacierService.GetJobOutput(svc, jobID)
	defer jobOutputOutput.Body.Close()
	var response, err = ioutil.ReadAll(jobOutputOutput.Body)
	ShowError(err)
	retrievedFile, err := ioutil.TempFile("/tmp", "movie.*.zip")
	ShowError(err)
	glacierService.Unzip(retrievedFile.Name())
	err = ioutil.WriteFile(retrievedFile.Name(), response, 0644)
	glacierService.Cleanup([]string{retrievedFile.Name()})
	ShowError(err)
}

type InventoryRetrieve struct {
	InventoryDate string                    `json:"InventoryDate"`
	VaultARN      string                    `json:"VaultARN"`
	ArchiveList   []models.InventoryArchive `json:"ArchiveList"`
}
