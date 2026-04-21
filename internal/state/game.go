package state

import (
	"sort"
	"time"

	"github.com/timkelleher/spaceui/internal/api"
)

const (
	DATA_SHIPS_LIST = "ships_list"
)

// /////////////////////////////////////
// Loading & Last Updated
// /////////////////////////////////////
var (
	loading     map[string]bool
	lastUpdated map[string]time.Time
	globalError string
)

// Newly proposed system:
// If currently loading, return true
// If not currently loading, but data has been update before, return false
// If not currently loading, but data has never been updated before, return true and trigger update
func Loading(id string) bool {
	loading, ok := loading[id]
	if !ok {
		return false
	}
	return loading
}

func LastUpdated(id string) time.Time {
	lastUpdated, ok := lastUpdated[id]
	if !ok {
		return time.Time{}
	}
	return lastUpdated
}

func GlobalError() string {
	return globalError
}

// /////////////////////////////////////
// Selected Ship
// /////////////////////////////////////
var selectedShip *api.Ship

func HasSelectedShip() bool {
	return SelectedShip() != nil
}

func SelectedShip() *api.Ship {
	return selectedShip
}

func SetSelectedShip(s *api.Ship) {
	selectedShip = s
}

func DeselectShip() {
	selectedShip = nil
}

// /////////////////////////////////////
// Active Ship
// /////////////////////////////////////
var activeShip *api.Ship

func HasActiveShip() bool {
	return ActiveShip() != nil
}

func ActiveShip() *api.Ship {
	return activeShip
}

func SetActiveShip(s *api.Ship) {
	activeShip = s
}

func DeactivateShip() {
	activeShip = nil
}

// /////////////////////////////////////
// Selected Waypoint
// /////////////////////////////////////
var selectedWaypoint *api.Waypoint

func HasSelectedWaypoint() bool {
	return SelectedWaypoint() != nil
}

func SelectedWaypoint() *api.Waypoint {
	return selectedWaypoint
}

func SetSelectedWaypoint(s *api.Waypoint) {
	selectedWaypoint = s
}

// /////////////////////////////////////
// Waypoint Filtering
// /////////////////////////////////////
var (
	waypointFilterName string
	waypointFilterType string
)

func WaypointFilterName() string {
	return waypointFilterName
}

func SetWaypointFilterName(name string) {
	waypointFilterName = name
}

func ResetWaypointFilters() {
	waypointFilterName = ""
	waypointFilterType = ""
}

func WaypointFilterType() string {
	return waypointFilterType
}

func SetWaypointFilterType(t string) {
	waypointFilterType = t
}

// /////////////////////////////////////
// Utils
// /////////////////////////////////////
func FilteredWaypoints() []api.Waypoint {
	ship := ActiveShip()
	waypoints := Waypoints(ship.Nav.SystemSymbol)

	var filtered []api.Waypoint
	for _, waypoint := range waypoints {
		if WaypointFilterType() == "type" && waypoint.Type == WaypointFilterName() {
			filtered = append(filtered, waypoint)
		} else if WaypointFilterType() == "trait" {
			filtered = WaypointsWithTrait(ship, WaypointFilterName())
		}
	}
	return filtered
}

type WaypointCount struct {
	Name  string
	Count int
}

func WaypointsByType(ship *api.Ship) []WaypointCount {
	waypoints := Waypoints(ship.Nav.SystemSymbol)

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

func findWaypoint(symbol string) *api.Waypoint {
	ship := ActiveShip()
	if ship == nil {
		return nil
	}

	waypoints := Waypoints(ship.Nav.SystemSymbol)
	for _, waypoint := range waypoints {
		if waypoint.Symbol == symbol {
			return &waypoint
		}
	}
	return nil
}

func WaypointsWithTrait(ship *api.Ship, filter string) []api.Waypoint {
	waypoints := Waypoints(ship.Nav.SystemSymbol)

	var filtered []api.Waypoint
	for _, waypoint := range waypoints {
		for _, trait := range waypoint.Traits {
			if trait.Name == filter {
				filtered = append(filtered, waypoint)
			}
		}
	}
	return filtered
}

func ShipsAtWaypoint(waypointSymbol string) []string {
	var atWaypoint []string

	ships := Ships(false)
	for _, ship := range ships {
		if ship.Nav.WaypointSymbol == waypointSymbol {
			atWaypoint = append(atWaypoint, ship.Symbol)
		}
	}
	return atWaypoint
}
