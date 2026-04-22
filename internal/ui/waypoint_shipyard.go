package ui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/state"
)

const (
	MENU_WAYPOINT_SHIPYARD = "waypoint_shipyard"
	PAGE_WAYPOINT_SHIPYARD = "waypoint_shipyard"
)

type WaypointShipyardMenu struct {
}

func (wsm WaypointShipyardMenu) ID() string {
	return MENU_WAYPOINT_SHIPYARD
}

func (wsm WaypointShipyardMenu) Menu() *tview.List {
	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			state.SetActivePage(PAGE_WAYPOINT_DETAIL)
			app.draw(true)
		})

	waypoint := state.SelectedWaypoint()
	availableShips := state.AvailableShips(waypoint)
	for _, ship := range availableShips.Ships {
		menu.AddItem(fmt.Sprintf("Buy %s", ship.Name), fmt.Sprintf("%d", ship.PurchasePrice), 0, func() {
			api.BuyShip(ship.Type, waypoint.Symbol)

			state.Queue(state.DATA_SHIPS)
			state.SetActivePage(PAGE_DASHBOARD)
			app.draw(true)

		})
	}

	return menu
}

type WaypointShipyardPage struct {
}

func (wsp WaypointShipyardPage) ID() string {
	return PAGE_WAYPOINT_SHIPYARD
}

func (wsp WaypointShipyardPage) Menu() Menu {
	return WaypointShipyardMenu{}
}

func (wsp WaypointShipyardPage) RequiredData() []string {
	return []string{state.DATA_WAYPOINTS}
}

func (wsp WaypointShipyardPage) Content() string {
	waypoint := state.SelectedWaypoint()
	if waypoint == nil {
		return "[red]Error: no active waypoint![-]"
	}

	shipyard := state.AvailableShips(waypoint)
	content := fmt.Sprintf("Available ship types for [yellow]%s[-] waypoint: [orange]%d[-]\n", waypoint.Symbol, len(shipyard.ShipTypes))
	for _, ship := range shipyard.ShipTypes {
		content += fmt.Sprintf("\t- [blue]%s[-]\n", ship.Type)
	}
	content += "\n"
	if len(shipyard.Ships) == 0 {
		content += "[red]Error: Shipyard has no ship data to view!  Transport a ship to this waypoint to view this data.[-]\n"
	}
	for _, ship := range shipyard.Ships {
		content += fmt.Sprintf("[orange]%s[-] | [blue]%s[-] | [green]%s[-] supply | Cost [orange]%d[-] \n", ship.Name, ship.Type, ship.Supply, ship.PurchasePrice)
		content += fmt.Sprintf("%s\n", ship.Description)
	}
	return content
}
