package library

import (
	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/eschizoid/flixctl/aws/glacier"
	"github.com/spf13/cobra"
)

var RefreshLibraryCmd = &cobra.Command{
	Use:   "refresh",
	Short: "To Refresh Media Library",
	Long:  "to refresh metadata media library.",
	Run: func(cmd *cobra.Command, args []string) {
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		ec2.CreateVault(svc)

		// 1. curl 'https://34-227-236-211.2b551790fd364f029ab0eb0a399e13d7.plex.direct:32400/applications/ExportTools/backgroundScan?title=Movies&random=2.48&key=3&sectiontype=movie&X-Plex-Product=Plex%20Web&X-Plex-Version=3.73.2&X-Plex-Client-Identifier=n0ivmpcq9s9f3q22gzrj5qwt&X-Plex-Platform=Chrome&X-Plex-Platform-Version=70.0&X-Plex-Sync-Version=2&X-Plex-Device=OSX&X-Plex-Device-Name=Chrome&X-Plex-Device-Screen-Resolution=1599x899%2C1680x1050&X-Plex-Token=8gQZ8gBpoyVNuuT7wJSt&X-Plex-Language=en' -H 'Accept: application/xml' -H 'Referer: https://app.plex.tv/' -H 'Origin: https://app.plex.tv' -H 'Accept-Language: en' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36' --compressed
		// 2. csv to struct https://github.com/gocarina/gocsv
		// 3. save structs to bolt db (add new column with glacier location)
	},
}
