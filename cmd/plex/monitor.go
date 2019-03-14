package plex

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/eschizoid/flixctl/models"
	"github.com/eschizoid/flixctl/worker"
	"github.com/jrudio/go-plex-client"
	"github.com/spf13/cobra"
)

var MonitorPlexCmd = &cobra.Command{
	Use:   "monitor",
	Short: "To Monitor Plex Sessions",
	Long:  "to monitor plex sessions and shut it down if no activity.",
	Run: func(cmd *cobra.Command, args []string) {
		plexClient, err := plex.New(fmt.Sprintf("https://%s:32400", os.Getenv("FLIXCTL_HOST")), os.Getenv("PLEX_TOKEN"))
		ShowError(err)
		m := make(map[string]interface{})
		sessions, _ := plexClient.GetSessions()
		now := getTime()
		if sessions.MediaContainer.Size > 0 {
			err = models.Database.SaveLastActiveSession(now)
			ShowError(err)
			m["plex_status"] = "active"
			m["last_activity"] = fmt.Sprintf("%d/%d/%d %d:%d", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute())
		} else {
			var lastActiveTime time.Time
			lastActiveTime, err = models.Database.GetLastActiveSession()
			ShowError(err)
			duration := time.Since(lastActiveTime).Minutes()
			//inactiveTime, _ := strconv.Atoi(maxInactiveTime)
			if int(duration) >= 30 {
				m["plex_status"] = "stopping"
				m["last_activity"] = lastActiveTime
				asyncShutdown()
			} else {
				m["plex_status"] = "running"
				m["last_activity"] = fmt.Sprintf("%d/%d/%d %d:%d", lastActiveTime.Day(), lastActiveTime.Month(), lastActiveTime.Year(), lastActiveTime.Hour(), lastActiveTime.Minute())
			}
		}
		jsonString, _ := json.Marshal(m)
		fmt.Println(string(jsonString))
	},
}

func asyncShutdown() {
	stopTask := func() interface{} {
		Stop()
		return "done shutting down plex!"
	}
	tasks := []worker.TaskFunction{stopTask}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	worker.PerformTasks(ctx, tasks)
}

func getTime() time.Time {
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		fmt.Println(err)
	}
	return time.Now().In(location)
}
