package ui

import (
	"sort"

	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/state"
)

type WaypointCount struct {
	Name  string
	Count int
}

func waypointsByType(ship *api.Ship) []WaypointCount {
	waypoints := state.Waypoints(ship.Nav.SystemSymbol)

	waypointsByType := make(map[string]int)
	for _, waypoint := range waypoints {
		waypointsByType[waypoint.Type]++
	}

	waypointNames := make([]WaypointCount, 0)
	for waypointType, count := range waypointsByType {
		waypointNames = append(waypointNames, WaypointCount{Name: waypointType, Count: count})
	}

	sort.Slice(waypointNames, func(i, j int) bool {
		return waypointNames[i].Name < waypointNames[j].Name
	})
	return waypointNames
}

func waypointsWithTrait(ship *api.Ship, desired string) []api.Waypoint {
	waypoints := state.Waypoints(ship.Nav.SystemSymbol)

	var filtered []api.Waypoint
	for _, waypoint := range waypoints {
		for _, trait := range waypoint.Traits {
			if trait.Name == desired {
				filtered = append(filtered, waypoint)
			}
		}
	}
	return filtered
}
