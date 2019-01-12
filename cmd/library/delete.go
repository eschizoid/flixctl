package library

import (
	"encoding/json"
	"fmt"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/spf13/cobra"
)

var DeleteArchiveLibraryCmd = &cobra.Command{
	Use:   "delete",
	Short: "To Delete Archives From Library",
	Long:  "to delete a show or a movie from the library.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		deleteArchiveOutput := glacierService.DeleteArchive(svc, archiveID)
		err := libraryService.DeleteteGlacierInventoryArchive(archiveID)
		ShowError(err)
		json, _ := json.Marshal(deleteArchiveOutput)
		fmt.Println(string(json))
	},
}
