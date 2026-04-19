package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/timkelleher/spaceui/internal/state"
)

func (a *App) dashboardContent() string {
	agent := state.Agent(false)

	content := "----- Best Space Traders App Ever: Dashboard -----\n"
	content += fmt.Sprintf("[green]Current Time:[-]\t %s\n", a.FormattedTime(time.Now()))
	content += fmt.Sprintf("[red]Agent Name:[-]\t\t %s\n", agent.Symbol)
	if a.ShipState.HasSelectedShip() {
		content += fmt.Sprintf("[orange]Selected Ship:[-]\t %s\n", a.ShipState.SelectedShipSymbol())
	}
	return content
}

func (a *App) agentContent() string {
	agent := state.Agent(false)

	content := "[yellow]Current Agent Information[-]\n" +
		fmt.Sprintf("[blue]Account ID:[-]\t\t%s\n", agent.AccountID) +
		fmt.Sprintf("[blue]Symbol:[-]\t\t\t%s\n", agent.Symbol) +
		fmt.Sprintf("[blue]Ship Count:[-]\t\t%d\n", agent.ShipCount) +
		fmt.Sprintf("[blue]Headquarters:[-]\t%s\n", agent.Headquarters) +
		fmt.Sprintf("[blue]Credits:[-]\t\t%d\n", agent.Credits)
	return content
}

func (a *App) contractsContent() string {
	contracts := state.Contracts(false)

	content := fmt.Sprintf("[yellow]Number of contracts:[-] %d\n", len(contracts))
	content += "----------\n"
	for i, contract := range contracts {
		content += fmt.Sprintf("[red]ID:[-]\t\t\t %s[-]\n", contract.ID)
		content += fmt.Sprintf("[blue]Type:[-]\t\t %s\n", contract.Type)
		content += fmt.Sprintf("[orange]State:[-]\t\t %s\n", contract.State())
		if !contract.Accepted {
			content += fmt.Sprintf("[green]Deadline to Accept:[-] %s\n", a.FormattedTime(contract.DeadlineToAccept))
		}
		if !contract.Fulfilled {
			content += fmt.Sprintf("[green]Expiration:[-]\t %s\n", a.FormattedTime(contract.Expiration))
		}
		content += fmt.Sprintf("[green]Deadline:[-]\t %s\n", a.FormattedTime(contract.Terms.Deadline))
		content += "----- Delivery Terms -----\n"
		for _, deliver := range contract.Terms.Deliver {
			content += fmt.Sprintf("[orange]%s[-]\t (%d/%d)\n", deliver.DestinationSymbol, deliver.UnitsFulfilled, deliver.UnitsRequired)
		}
		if i != len(contracts)-1 {
			content += "\n"
		}
	}
	return content
}

func (a *App) shipsContent() string {
	ships := state.Ships(false)

	content := fmt.Sprintf("[yellow]Number of ships owned:[-] %d\n", len(ships))
	if len(ships) == 0 {
		return content
	}

	ship := ships[a.ShipState.SelectedShipIndex()]

	content += fmt.Sprintf("Current Ship selected: [yellow]%s[-]\n", ships[a.ShipState.SelectedShipIndex()].Symbol)
	content += "----------\n"
	content += fmt.Sprintf("[red]Ship %s[-]\n", ship.Symbol)
	content += fmt.Sprintf("[orange]Frame:[-] %s | [orange]Reactor:[-] %s | [orange]Engine:[-] %s\n", ship.Frame.Name, ship.Reactor.Name, ship.Engine.Name)
	content += "----- Registration -----\n"
	content += fmt.Sprintf("[blue]Name:[-] %s | [blue]Faction Symbol:[-] %s | [blue]Role:[-] %s\n", ship.Registration.Name, ship.Registration.FactionSymbol, ship.Registration.Role)
	content += "----- Nav -----\n"
	content += fmt.Sprintf("[yellow]System:[-] %s | [yellow]Waypoint:[-] %s\n", ship.Nav.SystemSymbol, ship.Nav.WaypointSymbol)
	content += fmt.Sprintf("[yellow]Status:[-] %s | [yellow]Flight Mode:[-] %s\n", ship.Nav.Status, ship.Nav.FlightMode)
	content += "----- Cargo -----\n"
	content += fmt.Sprintf("[green]Capacity:[-] %d | [green]Units:[-] %d\n", ship.Cargo.Capacity, ship.Cargo.Units)
	content += "----- Fuel -----\n"
	content += fmt.Sprintf("[green]Capacity:[-] %d | [green]Current:[-] %d | [green]Consumed:[-] %d\n", ship.Fuel.Capacity, ship.Fuel.Current, ship.Fuel.Consumed.Amount)
	if a.ShipState.SelectedShipIndex() != len(ships)-1 {
		content += "\n"
	}

	return content
}

func (a *App) systemsContent() string {
	systems := state.Systems(false)

	content := fmt.Sprintf("[yellow]Number of systems scanned:[-] %d\n", len(systems))
	if len(systems) == 0 {
		return content
	}

	for _, system := range systems {
		content += fmt.Sprintf("[orange]%s[-] [blue]%s[-] [gold]%d,%d[-]\n", system.Symbol, system.Type, system.X, system.Y)
		waypointCounts := make(map[string]int)
		var waypointNames []string
		for _, waypoint := range system.Waypoints {
			waypointCounts[waypoint.Type]++
		}
		for waypointType, count := range waypointCounts {
			waypointNames = append(waypointNames, fmt.Sprintf("%d [purple]%s[-]", count, waypointType))
		}
		sort.Strings(waypointNames)
		if len(waypointNames) > 0 {
			content += strings.Join(waypointNames, " ") + "\n"
		}
	}
	return content
}

func (a *App) statusBarContent() string {
	currentShip := "[orange]no_ship_selected[-]"
	if a.ShipState.HasSelectedShip() {
		currentShip = fmt.Sprintf("[orange]%s[-]", a.ShipState.SelectedShipSymbol())
	}

	now := fmt.Sprintf("%s", time.Now().Format(time.RFC1123))
	agent := state.Agent(false)
	agentSymbol := fmt.Sprintf("[red]%s[-]", agent.Symbol)
	credits := fmt.Sprintf("[green]%d[-]", agent.Credits)
	content := now + "\t" + agentSymbol + "\t" + currentShip + "\t" + credits

	err := state.Get("global_error")
	if err != "" {
		content = fmt.Sprintf("[red]%s[-]", err)
	}
	return content
}
