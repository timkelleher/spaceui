package state

import (
	"math"
	"sort"
	"time"

	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/events"
)

var (
	loading     map[string]bool
	lastUpdated map[string]time.Time

	loc         *time.Location
	globalError string

	agent          api.Agent
	contracts      []api.Contract
	ships          []api.Ship
	systems        []api.System
	waypoints      []api.Waypoint
	availableShips api.AvailableShips
)

func Init() {
	loading = make(map[string]bool)
	lastUpdated = make(map[string]time.Time)

	go refreshAgentData()
	go refreshContractsData()
	go refreshShipData()
	go refreshSystemsData()
}

func refreshAgentData() {
	for {
		Agent(true)
		time.Sleep(20 * time.Second)
	}
}

func refreshContractsData() {
	for {
		Contracts(true)
		time.Sleep(45 * time.Second)
	}
}

func refreshShipData() {
	for {
		Ships(true)
		time.Sleep(30 * time.Second)
	}
}

func refreshSystemsData() {
	for {
		Systems(true)
		time.Sleep(30000 * time.Second)
	}
}

func GlobalError() string {
	return globalError
}

func Loc() *time.Location {
	return loc
}

func SetLoc(l *time.Location) {
	loc = l
}

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

func Agent(force bool) api.Agent {
	if _, ok := lastUpdated["agent"]; force || !ok {
		obj, _ := api.GetAgent()
		lastUpdated["agent"] = time.Now()
		agent = obj.Agent
	}
	return agent
}

func Contracts(force bool) []api.Contract {
	if _, ok := lastUpdated["contracts"]; force || !ok {
		obj, _ := api.GetContracts()
		lastUpdated["contracts"] = time.Now()
		contracts = obj.Contracts
	}
	return contracts
}

func Ships(force bool) []api.Ship {
	if _, ok := lastUpdated["ships"]; force || !ok {
		obj, _ := api.GetShips()
		lastUpdated["ships"] = time.Now()
		ships = obj.Ships
	}
	return ships
}

func Systems(force bool) []api.System {
	if check := lastUpdated["systems"]; check.IsZero() {
		systems = make([]api.System, 0)

		obj, res := api.GetSystems(1)
		systems = append(systems, obj.Systems...)

		for obj.Meta.Page < 5 {
			obj, res = api.GetSystems(obj.Meta.Page + 1)
			if res.Error() {
				globalError = "critical: failed to fetch systems"
				return nil
			}
			if len(obj.Systems) == 0 {
				break
			}
			systems = append(systems, obj.Systems...)
			time.Sleep(500 * time.Millisecond)
		}

		sort.Slice(systems, func(i, j int) bool {
			dist1 := math.Sqrt(math.Pow(float64(systems[i].X), 2.0) + math.Pow(float64(systems[i].Y), 2.0))
			dist2 := math.Sqrt(math.Pow(float64(systems[j].X), 2.0) + math.Pow(float64(systems[j].Y), 2.0))
			return dist1 < dist2
		})
		lastUpdated["systems"] = time.Now()
	}
	return systems
}

// Requires info about "other"/ui state.  Don't cache this?
func Waypoints(system string) []api.Waypoint {
	if check := lastUpdated["waypoints"]; check.IsZero() {
		waypoints = make([]api.Waypoint, 0)
		loading["waypoints"] = true

		obj, res := api.GetWaypoints(system, 1)
		lastUpdated["waypoints"] = time.Now()
		waypoints = obj.Waypoints

		for obj.Meta.Page*obj.Meta.Limit < obj.Meta.Total {
			obj, res = api.GetWaypoints(system, obj.Meta.Page+1)
			if res.Error() {
				globalError = "critical: failed to fetch waypoints"
				return nil
			}
			if len(obj.Waypoints) == 0 {
				break
			}

			waypoints = append(waypoints, obj.Waypoints...)
			time.Sleep(500 * time.Millisecond)
		}

		sort.Slice(waypoints, func(i, j int) bool {
			return waypoints[i].Type < waypoints[j].Type
		})
		lastUpdated["waypoints"] = time.Now()
		loading["waypoints"] = false
		events.NewEvent(events.EVENT_LOAD_WAYPOINTS_COMPLETE)
	}

	return waypoints
}

func WaypointsDataStatus() int {
	return len(waypoints)
}

func Reset(id string) {
	if _, ok := lastUpdated[id]; !ok {
		return
	}

	lastUpdated[id] = time.Time{}
	switch id {
	case "":
		waypoints = make([]api.Waypoint, 0)
	}
}

func AvailableShips(system, waypoint string) api.AvailableShips {
	if check := lastUpdated["available_ships"]; check.IsZero() {
		obj, _ := api.GetAvailableShips(system, waypoint)

		sort.Slice(waypoints, func(i, j int) bool {
			return waypoints[i].Type < waypoints[j].Type
		})

		availableShips = obj.AvailableShips
		lastUpdated["available_ships"] = time.Now()
	}

	return availableShips
}
