package state

import (
	"sort"
	"time"

	"github.com/timkelleher/spaceui/internal/api"
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

func AnyLoading(ids []string) bool {
	for _, id := range ids {
		if Loading(id) {
			return true
		} else if lastUpdated[id].IsZero() {
			Queue(id)
			return true
		}
	}
	return false
}

func LastUpdated(id string) time.Time {
	lastUpdated, ok := lastUpdated[id]
	if !ok {
		return time.Time{}
	}
	return lastUpdated
}

func Fresh(id string) bool {
	lastUpdated, ok := lastUpdated[id]
	if !ok || lastUpdated.IsZero() {
		return false
	}
	return true
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

var visitingSystems []string

func NewVisitingSystem(symbol string) {
	exists := false
	for _, sys := range visitingSystems {
		if sys == symbol {
			exists = true
		}
	}

	if !exists {
		visitingSystems = append(visitingSystems, symbol)
	}
}

func IsVisitingSystem(symbol string) bool {
	exists := false
	for _, sys := range visitingSystems {
		if sys == symbol {
			exists = true
		}
	}
	return exists
}

func ResetVisitingSystems() {
	visitingSystems = make([]string, 0)
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
	waypoints := Waypoints(false)

	var filtered []api.Waypoint
	for _, waypoint := range waypoints[ship.Nav.SystemSymbol] {
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
	waypoints := Waypoints(false)

	waypointsByType := make(map[string]int)
	for _, waypoint := range waypoints[ship.Nav.SystemSymbol] {
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

	waypoints := Waypoints(false)
	for _, waypoint := range waypoints[ship.Nav.SystemSymbol] {
		if waypoint.Symbol == symbol {
			return &waypoint
		}
	}
	return nil
}

func WaypointsWithTrait(ship *api.Ship, filter string) []api.Waypoint {
	waypoints := Waypoints(false)

	var filtered []api.Waypoint
	for _, waypoint := range waypoints[ship.Nav.SystemSymbol] {
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
