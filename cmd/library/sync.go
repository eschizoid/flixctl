package library

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
		instanceID := ec2Service.FetchInstanceID(svc, awsResourceTagNameValue)
		if ec2Status := ec2Service.Status(svc, instanceID); strings.EqualFold(ec2Status, ec2StatusRunning) {
			SyncMovieLibrary(0, awsSession)
			SyncMovieLibrary(1, awsSession)
		} else {
			m := make(map[string]string)
			m["plex_status"] = strings.ToLower(ec2StatusStopped)
			jsonString, _ := json.Marshal(m)
			fmt.Println("\n" + string(jsonString))
		}
	},
}

func SyncMovieLibrary(unwatched int, awsSession *sess.Session) {
	awsSession.Config.Endpoint = aws.String("http://dynamodb:8000")
	svc := dynamodb.New(awsSession)
	movies, _ := libraryService.GetLivePlexMovies(unwatched)
	for _, movie := range movies {
		err := libraryService.SavePlexMovie(movie, svc)
		ShowError(err)
	}
}
