package state

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/events"
)

const (
	DATA_AGENT                = "agent"
	DATA_CONTRACTS            = "contracts"
	DATA_SHIPS                = "ships"
	DATA_WAYPOINTS            = "waypoints"
	DATA_WAYPOINT_MARKETPLACE = "marketplace"
	DATA_WAYPOINT_SHIPYARD    = "shipyard"

	DATA_SHIPS_LIST = "ships_list"
)

var (
	mu  sync.Mutex
	ctx context.Context

	agent     api.Agent
	contracts []api.Contract
	ships     []api.Ship
	//systems   []api.System
	waypoints   map[string][]api.Waypoint
	marketplace map[string]api.Marketplace
	shipyard    map[string]api.Shipyard
)

func Init() {
	mu = sync.Mutex{}

	ctx = context.Background()
	loading = make(map[string]bool)
	lastUpdated = make(map[string]time.Time)

	marketplace = make(map[string]api.Marketplace)
	shipyard = make(map[string]api.Shipyard)

	go Poll()

	go refreshAgentData()
	go refreshContractsData()
	go refreshShipData()
	//go refreshSystemsData()
}

func refreshAgentData() {
	for {
		time.Sleep(30 * time.Second)
		MarkStale(DATA_AGENT)
	}
}

func refreshContractsData() {
	for {
		time.Sleep(120 * time.Second)
		MarkStale(DATA_CONTRACTS)
	}
}

func refreshShipData() {
	for {
		time.Sleep(20 * time.Second)
		MarkStale(DATA_SHIPS)
	}
}

/*
func refreshSystemsData() {
	for {
		Systems(true)
		time.Sleep(30000 * time.Second)
	}
}
*/

func dataLoad(datatype string, dataLoader func(context.Context)) {
	loading[datatype] = true
	dataLoader(ctx)
	loading[datatype] = false

	lastUpdated[datatype] = time.Now()
}

func Agent(force bool) api.Agent {
	if updated, ok := lastUpdated[DATA_AGENT]; !ok || force || updated.IsZero() {
		dataLoad(DATA_AGENT, func(context.Context) {
			obj, _ := api.GetAgent()
			agent = obj.Agent
		})
	}
	return agent
}

func Contracts(force bool) []api.Contract {
	if updated, ok := lastUpdated[DATA_CONTRACTS]; !ok || force || updated.IsZero() {
		dataLoad(DATA_CONTRACTS, func(context.Context) {
			obj, _ := api.GetContracts()
			contracts = obj.Contracts
		})
	}
	return contracts
}

func Ships(force bool) []api.Ship {
	if updated, ok := lastUpdated[DATA_SHIPS]; !ok || force || updated.IsZero() {
		dataLoad(DATA_SHIPS, func(context.Context) {
			obj, _ := api.GetShips()
			ships = obj.Ships

			ResetVisitingSystems()
			for _, ship := range ships {
				NewVisitingSystem(ship.Nav.SystemSymbol)
			}

			if HasSelectedShip() {
				for _, ship := range ships {
					if selectedShip.Symbol == ship.Symbol {
						selectedShip = &ship
						break
					}
				}
			}
		})
		events.NewEvent(events.EVENT_LOAD_SHIPS_COMPLETE)
	}
	return ships
}

/*
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
*/

// This doesn't go through queueDataLoad() as it is paginated...
// TODO: make a paginated
func Waypoints(force bool) map[string][]api.Waypoint {
	if updated, ok := lastUpdated[DATA_WAYPOINTS]; !ok || force || updated.IsZero() {
		loading[DATA_WAYPOINTS] = true

		waypoints = make(map[string][]api.Waypoint)
		for _, visitingSystem := range visitingSystems {
			obj, res := api.GetWaypoints(visitingSystem, 1)
			waypoints[visitingSystem] = append(waypoints[visitingSystem], obj.Waypoints...)
			for obj.Meta.Page*obj.Meta.Limit < obj.Meta.Total {
				obj, res = api.GetWaypoints(visitingSystem, obj.Meta.Page+1)
				if res.Error() {
					globalError = "critical: failed to fetch waypoints"
					return nil
				}
				if len(obj.Waypoints) == 0 {
					break
				}

				waypoints[visitingSystem] = append(waypoints[visitingSystem], obj.Waypoints...)
				time.Sleep(time.Second)
			}
			lastUpdated[DATA_WAYPOINTS] = time.Now()
		}

		// Sort all waypoints
		for _, visitingSystem := range visitingSystems {
			sort.Slice(waypoints[visitingSystem], func(i, j int) bool {
				return waypoints[visitingSystem][i].Type < waypoints[visitingSystem][j].Type
			})
		}

		lastUpdated[DATA_WAYPOINTS] = time.Now()
		loading[DATA_WAYPOINTS] = false
		events.NewEvent(events.EVENT_LOAD_WAYPOINTS_COMPLETE)
	}

	return waypoints
}

func Marketplace(waypoint *api.Waypoint) api.Marketplace {
	datatype := DATA_WAYPOINT_MARKETPLACE + "_" + waypoint.Symbol

	if updated, ok := lastUpdated[datatype]; !ok || updated.IsZero() {
		dataLoad(datatype, func(context.Context) {
			loading[DATA_WAYPOINT_MARKETPLACE] = true // Lock all marketplaces for now?
			obj, _ := api.GetMarketplace(*waypoint)
			marketplace[waypoint.Symbol] = obj.Marketplace
			loading[DATA_WAYPOINT_MARKETPLACE] = false
		})
	}
	return marketplace[waypoint.Symbol]
}

func AvailableShips(waypoint *api.Waypoint) api.Shipyard {
	datatype := DATA_WAYPOINT_SHIPYARD + "_" + waypoint.Symbol

	if updated, ok := lastUpdated[datatype]; !ok || updated.IsZero() {
		dataLoad(datatype, func(context.Context) {
			loading[DATA_WAYPOINT_SHIPYARD] = true // Lock all shipyards for now?
			obj, _ := api.GetShipyard(*waypoint)
			shipyard[waypoint.Symbol] = obj.Shipyard
			loading[DATA_WAYPOINT_SHIPYARD] = false
		})
	}
	return shipyard[waypoint.Symbol]
}

func RefreshStatus(datatype string) int {
	count := 0

	switch datatype {
	case DATA_WAYPOINTS:
		count := 0
		for _, systemWaypoints := range waypoints {
			count += len(systemWaypoints)
		}
	}

	return count
}
