package library

import (
	"fmt"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/eschizoid/flixctl/models"
	"github.com/spf13/cobra"
)

var ArchiveLibraryCmd = &cobra.Command{
	Use:   "archive",
	Short: "To Archive Movies Or Shows",
	Long:  "to archive movies or shows to the library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		zipFile := glacierService.Zip(fileName)
		zipFileName := zipFile.Name()
		fileChunks := glacierService.Chunk(zipFileName)
		svc := glacier.New(awsSession)
		initiateMultipartUploadOutput := glacierService.InitiateMultipartUploadInput(svc, zipFileName)
		fmt.Println(initiateMultipartUploadOutput.String())
		uploadID := *initiateMultipartUploadOutput.UploadId
		uploadMultipartPartOutputs := glacierService.UploadMultipartPartInput(svc, uploadID, fileChunks)
		fmt.Println(uploadMultipartPartOutputs)
		archiveCreationOutput := glacierService.CompleteMultipartUpload(svc, uploadID, zipFileName)
		fmt.Println(archiveCreationOutput)
		upload := new(models.Upload)
		upload.ArchiveCreationOutput = archiveCreationOutput
		err := libraryService.SaveGlacierMovie(*upload)
		ShowError(err)
		glacierService.Cleanup(append(fileChunks, zipFileName))
		close(shutdownCh)
	},
}
