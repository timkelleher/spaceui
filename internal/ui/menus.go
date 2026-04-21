package ui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/state"
)

func (us *UIState) SelectedMenu() string {
	return us.selectedMenu
}

// TODO: safely select navbars
func (us *UIState) SetSelectedNavbar(id string) {
	us.selectedMenu = id
}

func (us *UIState) MenuChanged() bool {
	return us.selectedMenu != us.DeterminedMenu()
}

func (us *UIState) DeterminedMenu() string {
	switch us.SelectedPanel() {
	case PANEL_CONTRACTS:
		return MENU_CONTRACTS
	case PANEL_SHIPS_LIST:
		return MENU_SHIPS_LIST
	case PANEL_SHIP_DETAIL:
		return MENU_SHIP_DETAIL
	case PANEL_WAYPOINTS_LIST:
		return MENU_WAYPOINTS_LIST
	case PANEL_WAYPOINT_DETAIL:
		return MENU_WAYPOINT_DETAIL
	case PANEL_WAYPOINT_AVAILABLE_SHIPS:
		return MENU_WAYPOINT_AVAILABLE_SHIPS
	default:
		return MENU_MAIN
	}
}

func (a *App) Menu() tview.Primitive {
	switch a.UIState.SelectedPanel() {
	case PANEL_CONTRACTS:
		return a.contractsMenu()
	case PANEL_SHIPS_LIST:
		return a.shipsListMenu()
	case PANEL_SHIP_DETAIL:
		return a.shipDetailMenu()
	case PANEL_WAYPOINTS_LIST:
		return a.waypointsListMenu()
	case PANEL_WAYPOINT_DETAIL:
		return a.waypointDetailMenu()
	case PANEL_WAYPOINT_AVAILABLE_SHIPS:
		return a.waypointAvailableShips()
	default:
		return a.mainMenu()
	}
}

func (a *App) UpdateMenu() {
	a.UIState.selectedMenu = a.UIState.DeterminedMenu()
	a.navbar = a.Menu()
}

func (a *App) mainMenu() *tview.List {
	menu := tview.NewList().
		AddItem("Dashboard", "", 'd', func() {
			a.SetSelectedPanel(PANEL_DASHBOARD)
		}).
		AddItem("Agent", "", 'a', func() {
			a.SetSelectedPanel(PANEL_AGENT)
		}).
		AddItem("Contracts", "", 'c', func() {
			a.SetSelectedPanel(PANEL_CONTRACTS)
		}).
		AddItem("Ships", "", 's', func() {
			a.SetSelectedPanel(PANEL_SHIPS_LIST)
		}).
		AddItem("Systems", "", 'y', func() {
			a.SetSelectedPanel(PANEL_SYSTEMS)
		})

	if a.GameState.HasActiveShip() {
		menu.AddItem("Waypoints", "", 'w', func() {
			a.SetSelectedPanel(PANEL_WAYPOINTS_LIST)
		})
	}

	menu.AddItem("Quit", "", 'q', func() {
		api.Close()
		a.ui.Stop()
	})

	return menu
}

func (a *App) contractsMenu() *tview.List {
	contracts := state.Contracts(false)

	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			a.SetSelectedPanel(PANEL_DASHBOARD)
			a.draw(true)
		})

	option := '1'
	for _, contract := range contracts {
		// only 9 options
		if option == '0' {
			break
		}
		if !contract.Accepted {
			menu.AddItem("Accept Contract", contract.ID, option, func() {
				api.AcceptContract(contract.ID)
				a.draw(true)
			})
			option++
		} else if !contract.Fulfilled {
			menu.AddItem("Fulfill Contract", contract.ID, option, func() {
				//api.FulfillContract(contract.ID)
				a.draw(true)
			})
			option++
		}
	}
	// TODO: auto select the command ship?
	if option != '0' {
		menu.AddItem("Negotiate New Contract", "", 'n', func() {
			api.NegotiateContract(a.GameState.ActiveShipSymbol())
			a.draw(true)
		})
	}
	return menu
}

func (a *App) shipsListMenu() *tview.List {
	ships := state.Ships(false)

	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			a.SetSelectedPanel(PANEL_DASHBOARD)
			a.draw(true)
		})

	activeShipSymbol := a.GameState.ActiveShipSymbol()

	for i, ship := range ships {
		desc := ""
		if activeShipSymbol == ship.Symbol {
			desc = "active"
		}
		menu.AddItem(ship.Symbol, desc, 0, func() {
			a.GameState.SelectShip(i)
			a.UIState.SetSelectedPanel(PANEL_SHIP_DETAIL)
			a.draw(true)
		})
	}

	return menu
}

func (a *App) shipDetailMenu() *tview.List {
	ships := state.Ships(false)

	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			a.SetSelectedPanel(PANEL_SHIPS_LIST)
			a.GameState.DeselectShip()
			a.draw(true)
		})

	selectedShip := ships[a.GameState.SelectedShipIndex()]
	if a.GameState.HasActiveShip() && a.GameState.ActiveShipSymbol() == selectedShip.Symbol {
		menu.AddItem("Deactivate Ship", "", 'd', func() {
			a.GameState.DeactivateShip()
			a.UIState.SetSelectedPanel(PANEL_SHIPS_LIST)
			a.draw(true)
		})
	} else {
		menu.AddItem("Activate Ship", "", 'a', func() {
			a.GameState.ActivateShip(a.GameState.SelectedShipIndex(), selectedShip.Symbol)
			a.UIState.SetSelectedPanel(PANEL_SHIP_DETAIL)
			a.draw(true)
		})
	}

	return menu
}

func (a *App) waypointsListMenu() *tview.List {
	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			a.SetSelectedPanel(PANEL_DASHBOARD)
			if a.GameState.waypointFilterName != "" {
				a.SetSelectedPanel(PANEL_WAYPOINTS_LIST)
			}

			a.GameState.waypointFilterType = ""
			a.GameState.waypointFilterName = ""
			a.draw(true)
		})

	ship := a.GameState.ActiveShip()
	if ship == nil {
		return menu
	}

	if !state.Loading("waypoints") && !state.LastUpdated("waypoints").IsZero() {
		waypoints := a.filteredWaypoints()
		for _, waypoint := range waypoints {
			menu.AddItem(waypoint.Symbol, "", 0, func() {
				a.SetSelectedPanel(PANEL_WAYPOINT_DETAIL)
				a.GameState.selectedWaypointSymbol = waypoint.Symbol
				a.draw(true)
			})
		}

		menu.AddItem("- Waypoints Filters -", "", 0, nil)

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
			waypointsByType := waypointsWithTrait(ship, desired)
			menu.AddItem(desired, fmt.Sprintf("%d", len(waypointsByType)), 0, func() {
				a.GameState.waypointFilterType = "trait"
				a.GameState.waypointFilterName = desired
				a.draw(true)
			})
		}

		// By Type
		waypointsByType := waypointsByType(ship)
		for _, waypointCount := range waypointsByType {
			menu.AddItem(waypointCount.Name, fmt.Sprintf("%d", waypointCount.Count), 0, func() {
				a.GameState.waypointFilterType = "type"
				a.GameState.waypointFilterName = waypointCount.Name
				a.draw(true)
			})
		}
	}

	return menu
}

func (a *App) waypointDetailMenu() *tview.List {
	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			a.SetSelectedPanel(PANEL_WAYPOINTS_LIST)
			a.draw(true)
		})

	waypoint := a.GameState.ActiveWaypoint()
	if waypoint.IsShipyard() {
		menu.AddItem("View Ships at Shipyard", "", 'v', func() {
			a.SetSelectedPanel(PANEL_WAYPOINT_AVAILABLE_SHIPS)
			a.draw(true)
		})
	}

	return menu
}

func (a *App) waypointAvailableShips() *tview.List {
	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			a.SetSelectedPanel(PANEL_WAYPOINT_DETAIL)
			a.draw(true)
		})

	waypoint := a.GameState.ActiveWaypoint()
	availableShips := state.AvailableShips(waypoint.SystemSymbol, waypoint.Symbol)
	for _, ship := range availableShips.ShipTypes {
		menu.AddItem(fmt.Sprintf("Buy %s", ship.Type), "", 'b', func() {
			// buy ship
		})
	}

	return menu
}
