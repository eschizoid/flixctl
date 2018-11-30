package library

import (
	"encoding/json"
	"fmt"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	"github.com/spf13/cobra"
)

var InitiateLibraryCmd = &cobra.Command{
	Use:   "initiate",
	Short: "To Initiate Library Catalogue Retrieval",
	Long:  "to initiate library catalogue retrieval.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		initiateJobOutput := glacierService.InitiateInventoryJob(svc)
		json, _ := json.Marshal(initiateJobOutput)
		fmt.Println(string(json))
	},
}
