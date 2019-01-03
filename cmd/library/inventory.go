package library

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	libraryService "github.com/eschizoid/flixctl/library"
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
		if sync, _ := strconv.ParseBool(enableLibrarySync); sync {
			err := libraryService.DeleteteGlacierInventoryArchives()
			ShowError(err)
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
		fmt.Println(string(jsonString))
		close(shutdownCh)
	},
}
