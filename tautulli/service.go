package tautulli

import (
	"fmt"

	"github.com/eschizoid/flixctl/slack/tautulli"
)

func CreateEvent(event ...string) {
	if err := tautulli.ForwardEvent(event[1]); err != nil {
		fmt.Println("Unable to forward message to slack")
	}
}
