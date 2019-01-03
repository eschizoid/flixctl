package library

import (
	"encoding/json"
	"fmt"
	"strings"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/eschizoid/flixctl/models"
	"github.com/spf13/cobra"
)

var CatalogueLibraryCmd = &cobra.Command{
	Use:   "catalogue",
	Short: "To Show Plex And Library Catalogue",
	Long:  "to show plex and library catalogue.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := ec2.New(awsSession, awsSession.Config)
		instanceID := ec2Service.FetchInstanceID(svc, "plex")
		var uploads []models.Upload
		var unwatchedMovies []models.Movie
		var err error
		if ec2Status := ec2Service.Status(svc, instanceID); strings.EqualFold(ec2Status, ec2StatusRunning) {
			unwatchedMovies, err = libraryService.GetLivePlexMovies(1)
			ShowError(err)
		} else {
			plexMovies, err := libraryService.GetCachedPlexMovies()
			unwatchedMovies = chooseUnwatchedMovies(plexMovies, func(unwatched int) bool {
				return unwatched == 1
			})
			ShowError(err)
			uploads, err = libraryService.GetGlacierMovies()
			ShowError(err)
			for _, upload := range uploads {
				glacierMovie := models.Movie{
					Unwatched: 0,
					Metadata:  upload.Metadata,
				}
				unwatchedMovies = append(unwatchedMovies, glacierMovie)
			}
		}
		json, _ := json.Marshal(unwatchedMovies)
		fmt.Println(string(json))
	},
}

func chooseUnwatchedMovies(movies []models.Movie, test func(int) bool) (unwatchedMovies []models.Movie) {
	for _, movie := range movies {
		if test(movie.Unwatched) {
			unwatchedMovies = append(unwatchedMovies, movie)
		}
	}
	return
}
