package library

import (
	"fmt"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	"github.com/spf13/cobra"
)

var ArchiveLibraryCmd = &cobra.Command{
	Use:   "archive",
	Short: "To Archive a File To Media Library",
	Long:  "to archive a movie or show to media library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		fileChunks := glacierService.Chunk(fileName)
		svc := glacier.New(awsSession)
		initiateMultipartUploadOutput := glacierService.InitiateMultipartUploadInput(svc)
		fmt.Println(initiateMultipartUploadOutput.String())
		uploadID := *initiateMultipartUploadOutput.UploadId
		uploadMultipartPartOutputs := glacierService.UploadMultipartPartInput(svc, uploadID, fileChunks)
		fmt.Println(uploadMultipartPartOutputs)
		glacierService.CompleteMultipartUpload(svc, uploadID, fileName)
		close(shutdownCh)
	},
}
