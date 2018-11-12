package library

import (
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	"github.com/spf13/cobra"
)

var InitiateLibraryCmd = &cobra.Command{
	Use:   "initiate",
	Short: "To Initiate a File Retrieval",
	Long:  "to initiate the retrieval of a movie or show from media library.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		glacierService.InitiateJob(svc, retrievalType, archiveID)
	},
}
