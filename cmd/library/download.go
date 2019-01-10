package library

import (
	"io/ioutil"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	"github.com/spf13/cobra"
)

var DownloadLibraryCmd = &cobra.Command{
	Use:   "download",
	Short: "To Download Movies Or Shows",
	Long:  "to download movies or shows from the library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		jobOutputOutput := glacierService.GetJobOutput(svc, jobID)
		defer jobOutputOutput.Body.Close()
		var response, err = ioutil.ReadAll(jobOutputOutput.Body)
		ShowError(err)
		err = ioutil.WriteFile(targetFile, response, 0644)
		glacierService.Unzip(targetFile)
		//glacierService.CleanupFiles([]string{targetFile}, "")
		ShowError(err)
		close(shutdownCh)
	},
}
