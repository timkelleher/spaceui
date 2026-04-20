package events

import (
	"fmt"

	"github.com/timkelleher/spaceui/internal/logger"
)

const (
	EVENT_LOAD_WAYPOINTS_COMPLETE = "load_waypoints_complete"
)

var events []string

func NewEvent(id string) {
	events = append(events, id)
	logger.Info(fmt.Sprintf("EVENT %s", id))
}

func GetEvents() []string {
	if len(events) == 0 {
		return nil
	}

	copiedEvents := make([]string, len(events))
	copy(copiedEvents, events)
	events = make([]string, 0)

	logger.Info(fmt.Sprintf("DELIVER %d events", len(copiedEvents)))
	return copiedEvents
}
