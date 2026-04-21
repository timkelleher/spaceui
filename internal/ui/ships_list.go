package ui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/state"
)

const (
	MENU_SHIPS_LIST = "ships_list"
	PAGE_SHIPS_LIST = "ships_list"
)

type ShipsListMenu struct {
}

func (slm ShipsListMenu) ID() string {
	return MENU_SHIPS_LIST
}

func (slm ShipsListMenu) Menu() *tview.List {
	ships := state.Ships(false)

	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			state.SetActivePage(PAGE_DASHBOARD)
			app.draw(true)
		})

	for _, ship := range ships {
		desc := ""
		if state.HasActiveShip() && state.ActiveShip().Symbol == ship.Symbol {
			desc = "active"
		}
		menu.AddItem(ship.Symbol, desc, 0, func() {
			state.SetSelectedShip(&ship)
			state.SetActivePage(PAGE_SHIP_DETAIL)
			app.draw(true)
		})
	}

	return menu
}

type ShipsListPage struct {
}

func (slp ShipsListPage) ID() string {
	return PAGE_SHIPS_LIST
}

func (slp ShipsListPage) Menu() Menu {
	return ShipsListMenu{}
}

func (slp ShipsListPage) Content() string {
	if state.Loading("waypoints") {
		return loadingContent(0, "ships")
	}

	ships := state.Ships(false)

	content := fmt.Sprintf("[yellow]Number of ships owned:[-] %d\n", len(ships))
	if len(ships) == 0 {
		return content
	}
	if state.HasActiveShip() {
		content += fmt.Sprintf("[yellow]Currently active ship:[-] [red]%s[-]\n", state.ActiveShip().Symbol)
	}

	content += "----------\n"
	for i, ship := range ships {
		content += fmt.Sprintf("%d) [red]%s[-] [blue]%s[-] [gold]%s[-] %s\n",
			i+1, ship.Symbol, ship.Registration.Role, ship.Nav.WaypointSymbol, ship.Nav.Status)
	}
	return content
}
