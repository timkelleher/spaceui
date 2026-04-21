package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/state"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	MENU_WAYPOINTS_LIST = "waypoints_list"
	PAGE_WAYPOINTS_LIST = "waypoints_list"
)

type WaypointsListMenu struct {
}

func (wlm WaypointsListMenu) ID() string {
	return MENU_WAYPOINTS_LIST
}

func (wlm WaypointsListMenu) Menu() *tview.List {
	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			state.SetActivePage(PAGE_DASHBOARD)
			if state.WaypointFilterName() != "" {
				state.SetActivePage(PAGE_WAYPOINTS_LIST)
			}

			state.ResetWaypointFilters()
			app.draw(true)
		})

	ship := state.ActiveShip()
	if ship == nil {
		return menu
	}

	if !state.Loading("waypoints") && !state.LastUpdated("waypoints").IsZero() {
		waypoints := state.FilteredWaypoints()
		if len(waypoints) > 0 {
			menu.AddItem("Waypoints", "", 0, nil)
		}
		for _, waypoint := range waypoints {
			shipsAtWaypoint := state.ShipsAtWaypoint(waypoint.Symbol)
			atWaypoint := ""
			if len(shipsAtWaypoint) > 1 {
				atWaypoint = "   Ship at Waypoint"
			} else if len(shipsAtWaypoint) > 0 {
				atWaypoint = "   Ship at Waypoint"
			}
			menu.AddItem(fmt.Sprintf(" - %s", waypoint.Symbol), atWaypoint, 0, func() {
				state.SetActivePage(PAGE_WAYPOINT_DETAIL)
				state.SetSelectedWaypoint(&waypoint)
				app.draw(true)
			})
		}

		menu.AddItem("Filters", "", 0, nil)

		// TODO: rewrite this to combine waypointsWithTrait & waypointsByType somehow

		// Remove Filter
		//if a.GameState.waypointFilterType != "" {
		//	menu.AddItem("Remove Filter", "", 'r', func() {
		//		a.GameState.waypointFilterType = ""
		//		a.GameState.waypointFilterName = ""
		//		a.draw(true)
		//	})
		//}

		// By Trait
		desiredTraits := []string{"Shipyard"}
		for _, desired := range desiredTraits {
			waypointsByType := state.WaypointsWithTrait(ship, desired)
			menu.AddItem(fmt.Sprintf("Trait: %s", desired), fmt.Sprintf("%d Matches", len(waypointsByType)), 0, func() {
				state.SetWaypointFilterType("trait")
				state.SetWaypointFilterName(desired)
				app.draw(true)
			})
		}

		caser := cases.Title(language.English)

		// By Type
		waypointsByType := state.WaypointsByType(ship)
		for _, waypointCount := range waypointsByType {
			name := strings.ReplaceAll(waypointCount.Name, "_", " ")
			menu.AddItem(fmt.Sprintf("Type: %s", caser.String(name)), fmt.Sprintf("%d Matches", waypointCount.Count), 0, func() {
				state.SetWaypointFilterType("type")
				state.SetWaypointFilterName(waypointCount.Name)
				app.draw(true)
			})
		}
	}

	return menu
}

type WaypointsListPage struct {
	app *App
}

func (wlp WaypointsListPage) ID() string {
	return PAGE_WAYPOINTS_LIST
}

func (wlp WaypointsListPage) Menu() Menu {
	return WaypointsListMenu{}
}

func (wlp WaypointsListPage) Content() string {
	ship := state.ActiveShip()
	if ship == nil {
		return "[red]Error: no active ship![-]"
	}

	if !state.Loading("waypoints") && !state.LastUpdated("waypoints").IsZero() {
		var content strings.Builder
		waypoints := state.Waypoints(ship.Nav.SystemSymbol)

		content.WriteString(fmt.Sprintf("[yellow]Number of waypoints in %s system:[-] %d\n", ship.Nav.SystemSymbol, len(waypoints)))

		if state.WaypointFilterName() != "" {
			filtered := state.FilteredWaypoints()
			content.WriteString(fmt.Sprintf("[purple]Applying `%s` filter[-]: %d\n", state.WaypointFilterName(), len(filtered)))

			for _, waypoint := range filtered {
				shipsAtWaypoint := state.ShipsAtWaypoint(waypoint.Symbol)
				atWaypoint := strings.Join(shipsAtWaypoint, ", ")

				content.WriteString(fmt.Sprintf("[orange]%s[-] [blue]%s[-] [orange]%s[-]\n", waypoint.Symbol, waypoint.Type, atWaypoint))
				traits := waypoint.AllTraits()
				if len(traits) > 0 {
					content.WriteString(fmt.Sprintf("\t%s\n", strings.Join(traits, ", ")))
				}
			}
			return content.String()
		}

		waypointsByType := state.WaypointsByType(ship)
		for _, waypointCount := range waypointsByType {
			content.WriteString(fmt.Sprintf("[blue]%s[-] %d[-]\n", waypointCount.Name, waypointCount.Count))
		}

		return content.String()
	}

	if state.Loading("waypoints") {
		return loadingContent(state.WaypointsDataStatus(), "waypoints")
	}

	// TODO: put this funcitonality into state pkg in a more explicit manner
	// like "RefreshData()"
	go state.Waypoints(ship.Nav.SystemSymbol)

	return loadingContent(state.WaypointsDataStatus(), "waypoints")
}
