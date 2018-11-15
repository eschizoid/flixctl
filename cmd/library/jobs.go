package library

import (
	"encoding/json"
	"fmt"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	"github.com/spf13/cobra"
)

var JobsLibraryCmd = &cobra.Command{
	Use:   "jobs",
	Short: "To List Media Library Jobs",
	Long:  "to list media library jobs.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		jobList := glacierService.ListJobs(svc)
		json, _ := json.Marshal(jobList)
		fmt.Println(string(json))
	},
}
