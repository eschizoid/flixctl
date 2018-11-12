package library

import (
	"io"
	"io/ioutil"

	sess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glacier"
	glacierService "github.com/eschizoid/flixctl/aws/glacier"
	"github.com/spf13/cobra"
)

const maxFileChunkSize = 1024 * 1024 * 4 // 4MB

var RetrieveLibraryCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "To Retrieve a File From Media Library",
	Long:  "to retrieve a movie or show from media library.",
	Run: func(cmd *cobra.Command, args []string) {
		shutdownCh := make(chan struct{})
		go Indicator(shutdownCh)
		var awsSession = sess.Must(sess.NewSessionWithOptions(sess.Options{
			SharedConfigState: sess.SharedConfigEnable,
		}))
		svc := glacier.New(awsSession)
		getJobOutputOutput := glacierService.GetJobOutput(svc, "")
		defer getJobOutputOutput.Body.Close()
		part, err := ioutil.ReadAll(io.LimitReader(getJobOutputOutput.Body, maxFileChunkSize))
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(fileName, part, 0644)
		if err != nil {
			panic(err)
		}
		close(shutdownCh)
	},
}
