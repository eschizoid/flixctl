package library

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/eschizoid/flixctl/models"
	slackService "github.com/eschizoid/flixctl/slack/library"
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
		svcEc2 := ec2.New(awsSession, awsSession.Config)
		awsSession.Config.Endpoint = aws.String(os.Getenv("DYNAMODB_ENDPOINT"))
		svcDynamo := dynamodb.New(awsSession)
		switch archiveFilter {
		case "all": //nolint:goconst
			cachedMovies, err := libraryService.GetCachedPlexMovies(svcDynamo)
			ShowError(err)
			var libraryMovies []models.Movie
			uploads, err := libraryService.GetGlacierMovies(svcDynamo)
			ShowError(err)
			for _, upload := range uploads {
				glacierMovie := models.Movie{
					Unwatched: 0,
					Metadata:  upload.Metadata,
				}
				libraryMovies = append(libraryMovies, glacierMovie)
			}
			allMovies := append(cachedMovies, libraryMovies...)
			if notify, _ := strconv.ParseBool(slackNotification); notify {
				slackService.SendCatalogue(allMovies, slackIncomingHookURL)
			}
			json, _ := json.Marshal(allMovies)
			fmt.Println(string(json))
		case "archived":
			var archivedMovies []models.Movie
			uploads, err := libraryService.GetGlacierMovies(svcDynamo)
			ShowError(err)
			for _, upload := range uploads {
				glacierMovie := models.Movie{
					Unwatched: 0,
					Metadata:  upload.Metadata,
				}
				archivedMovies = append(archivedMovies, glacierMovie)
			}
			if notify, _ := strconv.ParseBool(slackNotification); notify {
				slackService.SendCatalogue(archivedMovies, slackIncomingHookURL)
			}
			json, _ := json.Marshal(archivedMovies)
			fmt.Println(string(json))
		case "live":
			instanceID := ec2Service.FetchInstanceID(svcEc2, awsResourceTagNameValue)
			if ec2Status := ec2Service.Status(svcEc2, instanceID); strings.EqualFold(ec2Status, ec2StatusRunning) {
				liveMovies, err := libraryService.GetLivePlexMovies(1)
				ShowError(err)
				if notify, _ := strconv.ParseBool(slackNotification); notify {
					slackService.SendCatalogue(liveMovies, slackIncomingHookURL)
				}
				json, _ := json.Marshal(liveMovies)
				fmt.Println(string(json))
			} else {
				m := make(map[string]string)
				m["plex_status"] = strings.ToLower(ec2StatusStopped)
				jsonString, _ := json.Marshal(m)
				fmt.Println("\n" + string(jsonString))
			}
		case "watched":
			var watchedMovies []models.Movie
			plexMovies, err := libraryService.GetCachedPlexMovies(svcDynamo)
			ShowError(err)
			watchedMovies = filterMovies(plexMovies, func(unwatched int) bool {
				return unwatched == 0
			})
			if notify, _ := strconv.ParseBool(slackNotification); notify {
				slackService.SendCatalogue(watchedMovies, slackIncomingHookURL)
			}
			json, _ := json.Marshal(watchedMovies)
			fmt.Println(string(json))
		case "unwatched":
			var unwatchedMovies []models.Movie
			plexMovies, err := libraryService.GetCachedPlexMovies(svcDynamo)
			ShowError(err)
			unwatchedMovies = filterMovies(plexMovies, func(unwatched int) bool {
				return unwatched == 1
			})
			if notify, _ := strconv.ParseBool(slackNotification); notify {
				slackService.SendCatalogue(unwatchedMovies, slackIncomingHookURL)
			}
			json, _ := json.Marshal(unwatchedMovies)
			fmt.Println(string(json))
		default:
			plexMovies, err := libraryService.GetCachedPlexMovies(svcDynamo)
			ShowError(err)
			if notify, _ := strconv.ParseBool(slackNotification); notify {
				slackService.SendCatalogue(plexMovies, slackIncomingHookURL)
			}
			json, _ := json.Marshal(plexMovies)
			fmt.Println(string(json))
		}
	},
}

func filterMovies(movies []models.Movie, test func(int) bool) (unwatchedMovies []models.Movie) {
	for _, movie := range movies {
		if test(movie.Unwatched) {
			unwatchedMovies = append(unwatchedMovies, movie)
		}
	}
	return
}
