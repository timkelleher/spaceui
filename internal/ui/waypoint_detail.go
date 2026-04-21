package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
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

	waypoint := state.SelectedWaypoint()
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

func (wdp WaypointDetailPage) Content() string {
	ship := state.ActiveShip()
	if ship == nil {
		return "[red]Error: no active ship![-]"
	}

	waypoint := state.SelectedWaypoint()
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
