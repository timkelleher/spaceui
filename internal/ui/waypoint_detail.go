package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/state"
)

const (
	MENU_WAYPOINT_DETAIL = "waypoint_detail"
	PAGE_WAYPOINT_DETAIL = "waypoint_detail"
)

type WaypointDetailMenu struct {
}

func (wdm WaypointDetailMenu) ID() string {
	return MENU_WAYPOINTS_LIST
}

func (wdm WaypointDetailMenu) Menu() *tview.List {
	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			state.SetActivePage(PAGE_WAYPOINTS_LIST)
			app.draw(true)
		})

	ship := state.ActiveShip()
	waypoint := state.SelectedWaypoint()

	// Navigate to waypoint
	if ship.Nav.SystemSymbol == waypoint.SystemSymbol &&
		ship.Nav.WaypointSymbol != waypoint.Symbol &&
		ship.Nav.Status == "IN_ORBIT" {
		menu.AddItem(fmt.Sprintf("Navigate to %s", waypoint.Symbol), ship.Symbol, 'n', func() {
			api.NavigateShip(ship.Symbol, waypoint.Symbol)

			// Delay data refresh
			go func() {
				state.MarkStale(state.DATA_SHIPS)

				time.Sleep(2 * time.Second)
				state.Queue(state.DATA_SHIPS)
			}()

			state.SetActivePage(PAGE_SHIPS_LIST)
			app.draw(true)
		})
	}

	if waypoint.IsShipyard() {
		menu.AddItem("View Shipyard", "", 's', func() {
			state.SetActivePage(PAGE_WAYPOINT_SHIPYARD)
			app.draw(true)
		})
	}

	return menu
}

type WaypointDetailPage struct {
}

func (wdp WaypointDetailPage) ID() string {
	return PAGE_WAYPOINT_DETAIL
}

func (wdp WaypointDetailPage) Menu() Menu {
	return WaypointDetailMenu{}
}

func (wdp WaypointDetailPage) RequiredData() []string {
	return []string{state.DATA_WAYPOINTS}
}

func (wdp WaypointDetailPage) Content() string {
	ship := state.ActiveShip()
	if ship == nil {
		return "[red]Error: no active ship![-]"
	}

	waypoint := state.SelectedWaypoint()
	content := fmt.Sprintf("Waypoint: [yellow]%s[-]\n", waypoint.Symbol)
	content += fmt.Sprintf("\t[blue]Type:[-] %s\n", waypoint.Type)
	if len(waypoint.Traits) > 0 {
		content += fmt.Sprintf("\t[purple]Traits:[-] %s\n", strings.Join(waypoint.AllTraits(), ", "))
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
