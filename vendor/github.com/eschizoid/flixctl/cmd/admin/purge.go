package admin

import (
	"fmt"
	"os"

	"github.com/eschizoid/flixctl/worker"
	"github.com/spf13/cobra"
)

var RemotePurgeSlack = `sudo /bin/slack-cleaner --perform \
    --quiet \
    --token %s \
    --rate 2 \
    --message \
    --file \
    --channel %s \
    --bot \
    --user "*"`

var PurgeSlackCmd = &cobra.Command{
	Use:   "slack-purge",
	Short: "To purge slack messages",
	Long:  "To purge slack messages from all channels",
	Run: func(cmd *cobra.Command, args []string) {
		SlackPurge()
	},
}

func SlackPurge() {
	conn := GetSSHConnection()
	defer conn.Close()
	slackChannels := [4]string{"monitoring", "new-releases", "requests", "travis"}
	var tasks = make([]worker.TaskFunction, 0, len(slackChannels))
	for _, channel := range slackChannels {
		command := fmt.Sprintf(RemotePurgeSlack, os.Getenv("SLACK_LEGACY_TOKEN"), channel)
		message := fmt.Sprintf("Succesfully purged slack channel %s", channel)
		commandTask := func() interface{} {
			RunCommand(command, conn)
			return message
		}
		tasks = append(tasks, commandTask)
	}
	AsyncCommandExecution(tasks)
}
