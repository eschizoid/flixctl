package library

import (
	"fmt"
	"path/filepath"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/eschizoid/flixctl/models"
	"github.com/jrudio/go-plex-client" //nolint:goimports
	"github.com/spf13/cobra"
)

var UploadLibraryCmd = &cobra.Command{
	Use:   "upload",
	Short: "To Upload Movies Or Shows",
	Long:  "to upload movies or shows to the library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		movies, _ := libraryService.GetCachedPlexMovies()
		var upload models.Upload
		for _, movie := range movies {
			if movie.Unwatched == 0 {
				if sourceFile == "" {
					sourceFolder, err := filepath.Abs(filepath.Dir(movie.Metadata.Media[0].Part[0].File))
					ShowError(err)
					archiveCreationOutput := Archive(movie.Metadata.Title, sourceFolder)
					upload = models.Upload{
						ArchiveCreationOutput: *archiveCreationOutput,
						Metadata:              movie.Metadata,
					}
				} else {
					sourceFolder, err := filepath.Abs(filepath.Dir(sourceFile))
					ShowError(err)
					archiveCreationOutput := Archive(sourceFile, sourceFolder)
					upload = models.Upload{
						ArchiveCreationOutput: *archiveCreationOutput,
						Metadata: plex.Metadata{
							Title: sourceFile,
						},
					}
				}
				err := libraryService.SaveGlacierMovie(upload)
				ShowError(err)
				break
			}
		}
		close(shutdownCh)
	},
}

func Archive(fileDescription string, sourceFolder string) *glacier.ArchiveCreationOutput {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	zipFile := glacierService.Zip(sourceFolder)
	zipFileName := zipFile.Name()
	fileChunks := glacierService.Chunk(zipFileName)
	svc := glacier.New(awsSession)
	initiateMultipartUploadOutput := glacierService.InitiateMultipartUploadInput(svc, fileDescription)
	fmt.Println(initiateMultipartUploadOutput.String())
	uploadID := *initiateMultipartUploadOutput.UploadId
	uploadMultipartPartOutputs := glacierService.UploadMultipartPartInput(svc, uploadID, fileChunks)
	fmt.Println(uploadMultipartPartOutputs)
	archiveCreationOutput := glacierService.CompleteMultipartUpload(svc, uploadID, zipFileName)
	fmt.Println(archiveCreationOutput)
	glacierService.CleanupFiles(append(fileChunks, zipFileName), sourceFolder)
	return archiveCreationOutput
}
