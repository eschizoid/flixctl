package library

import (
	"encoding/json"
	"fmt"
	"strings"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	ec2Service "github.com/eschizoid/flixctl/aws/ec2"
	libraryService "github.com/eschizoid/flixctl/library"
	"github.com/jrudio/go-plex-client"
	"github.com/spf13/cobra"
)

var CatalogueLibraryCmd = &cobra.Command{
	Use:   "catalogue",
	Short: "To Get Movies And Shows Catalogue",
	Long:  "to get movies and shows catalogue from plex and library.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := ec2.New(awsSession, awsSession.Config)
		instanceID := ec2Service.FetchInstanceID(svc, "plex")
		var plexMovies []plex.Metadata
		var err error
		if ec2Status := ec2Service.Status(svc, instanceID); strings.EqualFold(ec2Status, ec2StatusRunning) {
			plexMovies, err = libraryService.GetLivePlexMovies("")
			ShowError(err)
		} else {
			plexMovies, err = libraryService.GetCachedPlexMovies()
			ShowError(err)
		}
		glacierMovies, err := libraryService.GetGlacierMovies()
		ShowError(err)
		if len(glacierMovies) > 0 {
			for _, glacierMovie := range glacierMovies {
				plexMovies = append(plexMovies, *glacierMovie.Metadata)
			}
		}
		json, _ := json.Marshal(plexMovies)
		fmt.Println(string(json))
	},
}
