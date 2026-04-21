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
	content += "\n"
	content += "----- Current Agent Information -----\n" +
		fmt.Sprintf("[blue]Symbol:[-]\t\t\t%s\n", agent.Symbol) +
		fmt.Sprintf("[blue]Ship Count:[-]\t\t%d\n", agent.ShipCount) +
		fmt.Sprintf("[blue]Headquarters:[-]\t%s\n", agent.Headquarters) +
		fmt.Sprintf("[blue]Credits:[-]\t\t%d\n", agent.Credits)
	content += "\n"
	if a.GameState.HasActiveShip() {
		ship := a.GameState.ActiveShip()
		content += "----- Current Ship Information -----\n"
		content += fmt.Sprintf("[orange]Active Ship:[-]\t %s\n", ship.Symbol)
		content += fmt.Sprintf("[orange]Status:[-]\t\t\t %s\n", ship.Nav.Status)
		content += fmt.Sprintf("[orange]Location:[-]\t %s\n", ship.Nav.WaypointSymbol)
		content += "\n"
	}
	if a.GameState.HasActiveWaypoint() {
		waypoint := a.GameState.ActiveWaypoint()
		content += "----- Current Waypoint Information -----\n"
		content += fmt.Sprintf("[red]Selected Waypoint:[-]\t %s\n", waypoint.Symbol)
		content += "\n"
	}

	return content
}

func (a *App) contractsContent() string {
	contracts := state.Contracts(false)

	content := fmt.Sprintf("[yellow]Number of contracts:[-] %d\n", len(contracts))
	content += "----------\n"

	if !a.GameState.HasActiveShip() {
		content += "\n" +
			"\t[yellow]Please activate a ship in order to negotiate a new contract.[-]\n" +
			"\n"
	}

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

func (a *App) shipsListContent() string {
	ships := state.Ships(false)

	content := fmt.Sprintf("[yellow]Number of ships owned:[-] %d\n", len(ships))
	if len(ships) == 0 {
		return content
	}
	if a.GameState.HasActiveShip() {
		content += fmt.Sprintf("[yellow]Currently active ship:[-] [red]%s[-]\n", a.GameState.ActiveShipSymbol())
	}

	content += "----------\n"

	for i, ship := range ships {
		content += fmt.Sprintf("%d) [red]%s[-] [blue]%s[-] [gold]%s[-] %s\n",
			i+1, ship.Symbol, ship.Registration.Role, ship.Nav.WaypointSymbol, ship.Nav.Status)
	}
	return content
}

func (a *App) shipDetailContent() string {
	ships := state.Ships(false)

	if len(ships) == 0 {
		return ""
	}

	ship := ships[a.GameState.SelectedShipIndex()]
	content := fmt.Sprintf("[red]Ship %s[-]\n", ship.Symbol)
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
	if a.GameState.SelectedShipIndex() != len(ships)-1 {
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

func (a *App) waypointsListContent() string {
	ship := a.GameState.ActiveShip()
	if ship == nil {
		return "[red]Error: no active ship![-]"
	}

	if !state.Loading("waypoints") && !state.LastUpdated("waypoints").IsZero() {
		var content strings.Builder
		waypoints := state.Waypoints(ship.Nav.SystemSymbol)

		content.WriteString(fmt.Sprintf("[yellow]Number of waypoints in %s system:[-] %d\n", ship.Nav.SystemSymbol, len(waypoints)))

		if a.GameState.waypointFilterName != "" {
			filtered := a.filteredWaypoints()
			content.WriteString(fmt.Sprintf("[purple]Applying `%s` filter[-]: %d\n", a.GameState.waypointFilterName, len(filtered)))

			for _, waypoint := range filtered {
				shipsAtWaypoint := a.GameState.ShipsAtWaypoint(waypoint.Symbol)
				atWaypoint := strings.Join(shipsAtWaypoint, ", ")

				content.WriteString(fmt.Sprintf("[orange]%s[-] [blue]%s[-] [orange]%s[-]\n", waypoint.Symbol, waypoint.Type, atWaypoint))
				traits := waypoint.AllTraits()
				if len(traits) > 0 {
					content.WriteString(fmt.Sprintf("\t%s\n", strings.Join(traits, ", ")))
				}
			}
			return content.String()
		}

		waypointsByType := waypointsByType(ship)
		for _, waypointCount := range waypointsByType {
			content.WriteString(fmt.Sprintf("[blue]%s[-] %d[-]\n", waypointCount.Name, waypointCount.Count))
		}

		return content.String()
	}

	if state.Loading("waypoints") {
		return loadingContent(state.WaypointsDataStatus())
	}

	go state.Waypoints(ship.Nav.SystemSymbol)
	return loadingContent(state.WaypointsDataStatus())
}

func (a *App) waypointDetailContent() string {
	ship := a.GameState.ActiveShip()
	if ship == nil {
		return "[red]Error: no active ship![-]"
	}

	waypoint := a.GameState.ActiveWaypoint()
	content := fmt.Sprintf("Waypoint: [yellow]%s[-] [blue]%s[-]\n", ship.Nav.SystemSymbol, waypoint.Type)
	if len(waypoint.Traits) > 0 {
		content += fmt.Sprintf("\t[blue]Traits:[-] %s\n", strings.Join(waypoint.AllTraits(), ", "))
	}
	content += fmt.Sprintf("\t[orange]Faction:[-] %s\n", waypoint.Faction.Symbol)

	// Ships at this waypoint
	ships := state.Ships(false)
	content += "\n" + fmt.Sprintf("Ship(s) at waypoint [yellow]%s[-]:\n", waypoint.Symbol)
	for _, ship := range ships {
		if ship.Nav.WaypointSymbol == waypoint.Symbol {
			content += fmt.Sprintf("\t[orange]%s[-]\n", ship.Symbol)
		}
	}
	if len(ships) == 0 {
		content += "\tNone\n"
	}

	return content
}

func (a *App) waypointAvailableShipsContent() string {
	waypoint := a.GameState.ActiveWaypoint()
	if waypoint == nil {
		return "[red]Error: no active waypoint![-]"
	}

	shipyard := state.AvailableShips(waypoint.SystemSymbol, waypoint.Symbol)
	content := fmt.Sprintf("Available ship types for [yellow]%s[-] waypoint: [orange]%d[-]\n", waypoint.Symbol, len(shipyard.ShipTypes))
	for _, ship := range shipyard.ShipTypes {
		content += fmt.Sprintf("\t- [blue]%s[-]\n", ship.Type)
	}
	content += "\n"
	if len(shipyard.Ships) == 0 {

	}
	for _, ship := range shipyard.Ships {
		content += fmt.Sprintf("[orange]%s[-] | [blue]%s[-] | [green]%s[-] supply | Cost [orange]%d[-] \n", ship.Name, ship.Type, ship.Supply, ship.PurchasePrice)
		content += fmt.Sprintf("%s\n", ship.Description)
	}
	return content
}

func (a *App) statusBarContent() string {
	currentShip := ""
	if a.GameState.HasActiveShip() {
		currentShip = fmt.Sprintf("[orange]%s[-]", a.GameState.ActiveShipSymbol())
	}

	now := fmt.Sprintf("%s", time.Now().Format(time.RFC1123))
	agent := state.Agent(false)
	agentSymbol := fmt.Sprintf("[red]%s[-]", agent.Symbol)
	credits := fmt.Sprintf("[green]%d[-]", agent.Credits)
	content := now + "\t" + agentSymbol + "\t" + currentShip + "\t" + credits

	err := state.GlobalError()
	if err != "" {
		content = fmt.Sprintf("[red]%s[-]", err)
	}
	return content
}

func loadingContent(num int) string {
	return fmt.Sprintf("Loading %d objects...\n", num)
}
