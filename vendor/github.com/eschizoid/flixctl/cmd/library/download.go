package library

import (
	"io"
	"os"

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
		from := glacierService.GetJobOutput(svc, jobID).Body
		defer from.Close()
		to, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE, 0666)
		ShowError(err)
		defer to.Close()
		_, err = io.Copy(to, from)
		ShowError(err)
		glacierService.Unzip(targetFile)
		glacierService.CleanupFiles([]string{targetFile}, "")
		ShowError(err)
		close(shutdownCh)
	},
}
