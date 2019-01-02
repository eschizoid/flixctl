package library

import (
	"fmt"
	"path/filepath"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/eschizoid/flixctl/models"
	"github.com/jrudio/go-plex-client"
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
		for _, movie := range movies {
			if movie.Unwatched == 0 {
				fmt.Println(movie.Metadata.Media[0].Part[0].File)
				//Archive(movie.Metadata)
			}
		}
		close(shutdownCh)
	},
}

func Archive(metadata plex.Metadata) {
	var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
		SharedConfigState: sess.SharedConfigEnable,
	}))
	sourceFolder, err := filepath.Abs(filepath.Dir(metadata.Media[0].Part[0].File))
	ShowError(err)
	zipFile := glacierService.Zip(sourceFolder)
	zipFileName := zipFile.Name()
	fileChunks := glacierService.Chunk(zipFileName)
	svc := glacier.New(awsSession)
	initiateMultipartUploadOutput := glacierService.InitiateMultipartUploadInput(svc, metadata)
	fmt.Println(initiateMultipartUploadOutput.String())
	uploadID := *initiateMultipartUploadOutput.UploadId
	uploadMultipartPartOutputs := glacierService.UploadMultipartPartInput(svc, uploadID, fileChunks)
	fmt.Println(uploadMultipartPartOutputs)
	archiveCreationOutput := glacierService.CompleteMultipartUpload(svc, uploadID, zipFileName)
	fmt.Println(archiveCreationOutput)
	upload := models.Upload{
		ArchiveCreationOutput: *archiveCreationOutput,
		Metadata:              metadata,
	}
	err = libraryService.SaveGlacierMovie(upload)
	ShowError(err)
	//glacierService.Cleanup(append(fileChunks, zipFileName))
}
