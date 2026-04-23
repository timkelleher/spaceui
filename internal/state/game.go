package state

import (
	"sort"

	"github.com/timkelleher/spaceui/internal/api"
)

// /////////////////////////////////////
// Error Handling
// /////////////////////////////////////
var globalError string

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

func IsSelectedActiveShip() bool {
	if !HasSelectedShip() || !HasActiveShip() {
		return false
	}

	return SelectedShip().Symbol == ActiveShip().Symbol
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
// Waypoint
// /////////////////////////////////////

func GetWaypoint(systemSymbol, waypointSymbol string) *api.Waypoint {
	waypoints, ok := waypoints[systemSymbol]
	if !ok {
		return nil
	}

	for _, waypoint := range waypoints {
		if waypoint.Symbol == waypointSymbol {
			return &waypoint
		}
	}
	return nil
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
