package ui

import (
	"fmt"
	"time"

	"github.com/timkelleher/spaceui/internal/state"
)

const PAGE_DASHBOARD = "dashboard"

type DashboardPage struct {
}

func (dp DashboardPage) ID() string {
	return PAGE_DASHBOARD
}

func (dp DashboardPage) Menu() Menu {
	return MainMenu{}
}

func (dp DashboardPage) RequiredData() []string {
	return []string{state.DATA_AGENT, state.DATA_CONTRACTS, state.DATA_SHIPS}
}

func (dp DashboardPage) Content() string {
	agent := state.Agent(false)

	content := "----- Best Space Traders App Ever: Dashboard -----\n"
	content += fmt.Sprintf("[green]Current Time:[-]\t %s\n", state.FormattedTime(time.Now()))
	content += "\n"
	content += "----- Current Agent Information -----\n" +
		fmt.Sprintf("[blue]Symbol:[-]\t\t\t%s\n", agent.Symbol) +
		fmt.Sprintf("[blue]Ship Count:[-]\t\t%d\n", agent.ShipCount) +
		fmt.Sprintf("[blue]Headquarters:[-]\t%s\n", agent.Headquarters) +
		fmt.Sprintf("[blue]Credits:[-]\t\t%d\n", agent.Credits)
	content += "\n"
	if state.HasActiveShip() {
		ship := state.ActiveShip()
		content += "----- Current Ship Information -----\n"
		content += fmt.Sprintf("[orange]Active Ship:[-]\t %s\n", ship.Symbol)
		content += fmt.Sprintf("[orange]Status:[-]\t\t\t %s\n", ship.Nav.Status)
		content += fmt.Sprintf("[orange]Location:[-]\t %s\n", ship.Nav.WaypointSymbol)
		content += "\n"
	}
	if state.HasSelectedWaypoint() {
		waypoint := state.SelectedWaypoint()
		content += "----- Current Waypoint Information -----\n"
		content += fmt.Sprintf("[red]Selected Waypoint:[-]\t %s\n", waypoint.Symbol)
		content += "\n"
	}

	return content
}
