package state

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/logger"
)

var (
	state     map[string]any
	agent     api.Agent
	contracts []api.Contract
	ships     []api.Ship
	systems   []api.System
)

func Init() {
	state = make(map[string]any)
	state["global_error"] = ""

	go refreshAgentData()
	go refreshContractsData()
	go refreshShipData()
	go refreshSystemData()
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

func refreshSystemData() {
	for {
		Systems(true)
		time.Sleep(3000 * time.Second)
	}
}

func Get(id string) any {
	if val, ok := state[id]; !ok {
		return "not found"
	} else {
		return val
	}
}

func Set(id string, val any) {
	state[id] = val
	logger.Info(fmt.Sprintf("STATE %s %s", id, val))
}

func Agent(force bool) api.Agent {
	if _, ok := state["last_fetched.agent"]; force || !ok {
		obj, _ := api.GetAgent()
		state["last_fetched.agent"] = time.Now()
		agent = obj.Agent
	}
	return agent
}

func Contracts(force bool) []api.Contract {
	if _, ok := state["last_fetched.contracts"]; force || !ok {
		obj, _ := api.GetContracts()
		state["last_fetched.contracts"] = time.Now()
		contracts = obj.Contracts
	}
	return contracts
}

func Ships(force bool) []api.Ship {
	if _, ok := state["last_fetched.ships"]; force || !ok {
		obj, _ := api.GetShips()
		state["last_fetched.ships"] = time.Now()
		ships = obj.Ships
	}
	return ships
}

func Systems(force bool) []api.System {
	if _, ok := state["last_fetched.systems"]; force || !ok {
		systems = make([]api.System, 0)

		obj, res := api.GetSystems(1)
		systems = append(systems, obj.Systems...)

		for obj.Meta.Page < 5 {
			obj, res = api.GetSystems(obj.Meta.Page + 1)
			if res.Error() {
				state["global_error"] = "critical: failed to fetch systems"
				return nil
			}
			if len(obj.Systems) == 0 {
				break
			}
			systems = append(systems, obj.Systems...)
			time.Sleep(time.Second)
		}
		state["last_fetched.systems"] = time.Now()
		systems = obj.Systems

		sort.Slice(systems, func(i, j int) bool {
			dist1 := math.Sqrt(math.Pow(float64(systems[i].X), 2.0) + math.Pow(float64(systems[i].Y), 2.0))
			dist2 := math.Sqrt(math.Pow(float64(systems[j].X), 2.0) + math.Pow(float64(systems[j].Y), 2.0))
			return dist1 < dist2
		})
	}
	return systems
}
