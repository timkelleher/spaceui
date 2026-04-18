package state

import (
	"strconv"
	"time"

	"github.com/timkelleher/spaceui/internal/api"
)

var state map[string]string
var ships []api.Ship

func Init() {
	state = make(map[string]string)
	state["globalError"] = ""

	go refreshAgentData()
	go refreshShipData()
}

func refreshAgentData() {
	for {
		agent, resp := api.GetAgent()
		if resp.Err != nil {
			state["globalError"] = resp.Err.Error()
			break
		}
		state["agent.account_id"] = agent.Agent.AccountID
		state["agent.credits"] = strconv.Itoa(agent.Agent.Credits)
		state["agent.ship_count"] = strconv.Itoa(agent.Agent.ShipCount)
		state["agent.starting_faction"] = agent.Agent.StartingFaction
		state["agent.symbol"] = agent.Agent.Symbol

		time.Sleep(30 * time.Second)
	}
}

func refreshShipData() {
	for {
		res, resp := api.GetShips()
		if resp.Err != nil {
			state["globalError"] = resp.Err.Error()
			break
		}

		ships = res.Ships
		time.Sleep(30 * time.Second)
	}
}

func Get(id string) string {
	if val, ok := state[id]; !ok {
		return "not found"
	} else {
		return val
	}
}

func Ships() []api.Ship {
	return ships
}
