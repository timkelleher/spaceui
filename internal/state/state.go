package state

import (
	"fmt"
	"time"

	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/logger"
)

var (
	state     map[string]any
	agent     *api.Agent
	contracts []api.Contract
	ships     []api.Ship
)

func Init() {
	state = make(map[string]any)
	state["global_error"] = ""

	go refreshAgentData()
	go refreshContractsData()
	go refreshShipData()
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

func Agent(force bool) *api.Agent {
	if _, ok := state["last_fetched.agent"]; force || !ok {
		res, resp := api.GetAgent()
		if resp.Err != nil {
			state["global_error"] = resp.Err.Error()
			return nil
		}
		state["last_fetched.agent"] = time.Now()
		agent = &res.Agent
	}
	return agent
}

func Contracts(force bool) []api.Contract {
	if _, ok := state["last_fetched.contracts"]; force || !ok {
		res, resp := api.GetContracts()
		if resp.Err != nil {
			state["global_error"] = resp.Err.Error()
			return nil
		}
		state["last_fetched.contracts"] = time.Now()
		contracts = res.Contracts
	}
	return contracts
}

func Ships(force bool) []api.Ship {
	if _, ok := state["last_fetched.ships"]; force || !ok {
		res, resp := api.GetShips()
		if resp.Err != nil {
			state["global_error"] = resp.Err.Error()
			return nil
		}
		state["last_fetched.ships"] = time.Now()
		ships = res.Ships
	}
	return ships
}
