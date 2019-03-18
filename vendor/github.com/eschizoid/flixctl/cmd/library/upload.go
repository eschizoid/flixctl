package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/eschizoid/flixctl/models"
	"github.com/spf13/cobra"
)

var UploadLibraryCmd = &cobra.Command{
	Use:   "upload",
	Short: "To Upload Movies Or Shows",
	Long:  "to upload movies or shows to the library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		awsSession.Config.Endpoint = aws.String(os.Getenv("DYNAMODB_ENDPOINT"))
		svc := dynamodb.New(awsSession)
		movies, _ := libraryService.GetCachedPlexMovies(svc)
		var upload models.Upload
		if batchMode, _ := strconv.ParseBool(enableBatchUpload); batchMode {
			if maxUploadItems != "" {
				index, err := strconv.Atoi(maxUploadItems)
				ShowError(err)
				if index <= len(movies) {
					movies = movies[:index]
				}
			}
			for _, movie := range movies {
				if movie.Unwatched == 0 {
					sourceFolder, err := filepath.Abs(filepath.Dir(movie.Metadata.Media[0].Part[0].File))
					ShowError(err)
					archiveCreationOutput := Archive(movie.Metadata.Title, sourceFolder)
					upload = models.Upload{
						ArchiveCreationOutput: *archiveCreationOutput,
						Metadata:              movie.Metadata,
						Title:                 movie.Title,
					}
				}
				err := libraryService.SaveGlacierMovie(upload, svc)
				ShowError(err)
			}
		}
		close(shutdownCh)
	},
}

func Archive(fileName string, sourceFolder string) *glacier.ArchiveCreationOutput {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	zipFileName := glacierService.Zip(sourceFolder)
	fileChunks := glacierService.Chunk(zipFileName)
	svc := glacier.New(awsSession)
	initiateMultipartUploadOutput := glacierService.InitiateMultipartUploadInput(svc, fileName)
	fmt.Println(initiateMultipartUploadOutput.String())
	uploadID := *initiateMultipartUploadOutput.UploadId
	uploadMultipartPartOutputs := glacierService.UploadMultipartPartInput(svc, uploadID, fileChunks)
	fmt.Println(uploadMultipartPartOutputs)
	archiveCreationOutput := glacierService.CompleteMultipartUpload(svc, uploadID, zipFileName)
	fmt.Println(archiveCreationOutput)
	glacierService.CleanupFiles(append(fileChunks, zipFileName), sourceFolder)
	return archiveCreationOutput
}
