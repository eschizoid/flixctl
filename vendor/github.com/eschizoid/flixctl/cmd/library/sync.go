package library

import (
	"strings"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/spf13/cobra"
)

var SyncLibraryCmd = &cobra.Command{
	Use:   "sync",
	Short: "To Sync Plex Watched Movies And Shows",
	Long:  "to sync plex watched movies and shows with internal flixctl db.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := ec2.New(awsSession, awsSession.Config)
		instanceID := ec2Service.FetchInstanceID(svc, "plex")
		if ec2Status := ec2Service.Status(svc, instanceID); strings.EqualFold(ec2Status, ec2StatusRunning) {
			movies, _ := libraryService.GetLivePlexMovies("?unwatched=1")
			for _, movie := range movies {
				err := libraryService.SavePlexMovie(movie)
				ShowError(err)
			}
		}
	},
}
