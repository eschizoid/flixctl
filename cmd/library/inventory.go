package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/eschizoid/flixctl/models"
	slackService "github.com/eschizoid/flixctl/slack/library"
	"github.com/spf13/cobra"
)

var InventoryLibraryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "To Show Library Inventory",
	Long:  "to show library inventory and if specified, sync with glacier.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svcGlacier := glacier.New(awsSession)
		awsSession.Config.Endpoint = aws.String("http://dynamodb:8000")
		svcDynamo := dynamodb.New(awsSession)
		if sync, _ := strconv.ParseBool(enableLibrarySync); sync && jobID != "" {
			err := libraryService.DeleteteGlacierInventoryArchives(svcDynamo)
			ShowError(err)
			jobOutputOutput := glacierService.GetJobOutput(svcGlacier, jobID)
			defer jobOutputOutput.Body.Close()
			response, err := ioutil.ReadAll(jobOutputOutput.Body)
			ShowError(err)
			var inventoryRetrieve = new(InventoryRetrieve)
			err = json.Unmarshal(response, &inventoryRetrieve)
			ShowError(err)
			for _, archive := range inventoryRetrieve.ArchiveList {
				err = libraryService.SaveGlacierInventoryArchive(archive, svcDynamo)
				ShowError(err)
			}
		} else {
			panic("job-id parameter should be provided")
		}
		glacierArchives, err := libraryService.GetGlacierInventoryArchives(svcDynamo)
		ShowError(err)
		jsonString, err := json.Marshal(glacierArchives)
		ShowError(err)
		if notify, _ := strconv.ParseBool(slackNotification); notify {
			slackService.SendInventory(glacierArchives, slackIncomingHookURL)
		}
		fmt.Println(string(jsonString))
		close(shutdownCh)
	},
}

type InventoryRetrieve struct {
	InventoryDate string                    `json:"InventoryDate"`
	VaultARN      string                    `json:"VaultARN"`
	ArchiveList   []models.InventoryArchive `json:"ArchiveList"`
}
