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
		if ec2Status := ec2Service.Status(svc, instanceID); strings.EqualFold(ec2Status, ec2StatusRunning) {
			liveMovies, err := libraryService.GetLivePlexMovies(1)
			ShowError(err)
			json, _ := json.Marshal(liveMovies)
			fmt.Println(string(json))
		} else {
			switch archiveFilter {
			case "all": //nolint:goconst
				cachedMovies, err := libraryService.GetCachedPlexMovies()
				ShowError(err)
				var libraryMovies []models.Movie
				uploads, err := libraryService.GetGlacierMovies()
				ShowError(err)
				for _, upload := range uploads {
					glacierMovie := models.Movie{
						Unwatched: 0,
						Metadata:  upload.Metadata,
					}
					libraryMovies = append(libraryMovies, glacierMovie)
				}
				json, _ := json.Marshal(append(cachedMovies, libraryMovies...))
				fmt.Println(string(json))
			case "archived":
				var archivedMovies []models.Movie
				uploads, err := libraryService.GetGlacierMovies()
				ShowError(err)
				for _, upload := range uploads {
					glacierMovie := models.Movie{
						Unwatched: 0,
						Metadata:  upload.Metadata,
					}
					archivedMovies = append(archivedMovies, glacierMovie)
				}
				json, _ := json.Marshal(archivedMovies)
				fmt.Println(string(json))
			case "live":
				cachedMovies, err := libraryService.GetCachedPlexMovies()
				ShowError(err)
				json, _ := json.Marshal(cachedMovies)
				fmt.Println(string(json))
			case "watched":
				var watchedMovies []models.Movie
				plexMovies, err := libraryService.GetCachedPlexMovies()
				ShowError(err)
				watchedMovies = chooseMovies(plexMovies, func(unwatched int) bool {
					return unwatched == 0
				})
				json, _ := json.Marshal(watchedMovies)
				fmt.Println(string(json))
			case "unwatched":
				var unwatchedMovies []models.Movie
				plexMovies, err := libraryService.GetCachedPlexMovies()
				ShowError(err)
				unwatchedMovies = chooseMovies(plexMovies, func(unwatched int) bool {
					return unwatched == 1
				})
				json, _ := json.Marshal(unwatchedMovies)
				fmt.Println(string(json))
			default:
				plexMovies, err := libraryService.GetCachedPlexMovies()
				ShowError(err)
				json, _ := json.Marshal(plexMovies)
				fmt.Println(string(json))
			}
		}
	},
}

func chooseMovies(movies []models.Movie, test func(int) bool) (unwatchedMovies []models.Movie) {
	for _, movie := range movies {
		if test(movie.Unwatched) {
			unwatchedMovies = append(unwatchedMovies, movie)
		}
	}
	return
}
