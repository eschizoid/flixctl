package admin

import (
	"fmt"

	"github.com/eschizoid/flixctl/worker"
	"github.com/spf13/cobra"
)

var RemoteRestartPlexServices = `sudo systemctl restart %s`

var RestartPlexServicesCmd = &cobra.Command{
	Use:   "restart-services",
	Short: "To Restart Plex Services",
	Long:  "to restart all plex related services",
	Run: func(cmd *cobra.Command, args []string) {
		RestartPlexServices()
	},
}

func RestartPlexServices() {
	conn := GetSSHConnection()
	defer conn.Close()
	services := [10]string{"httpd", "jackett", "nzbget", "ombi", "plexmediaserver", "radarr", "sonarr", "s3fs", "tautulli", "transmission-daemon"}
	var tasks = make([]worker.TaskFunction, 0, len(services))
	for _, service := range services {
		command := fmt.Sprintf(RemoteRestartPlexServices, service)
		message := fmt.Sprintf("Succesfully restarted service %s", service)
		commandTask := func() interface{} {
			RunCommand(command, conn)
			return message
		}
		tasks = append(tasks, commandTask)
	}
	AsyncCommandExecution(tasks)
}
