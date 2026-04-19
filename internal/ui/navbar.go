package ui

import (
	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/state"
)

func (us *UIState) SelectedNavbar() string {
	return us.selectedNavbar
}

// TODO: safely select navbars
func (us *UIState) SetSelectedNavbar(id string) {
	us.selectedNavbar = id
}

func (us *UIState) NavbarChanged() bool {
	return us.selectedNavbar != us.DeterminedNavbar()
}

func (us *UIState) DeterminedNavbar() string {
	switch us.SelectedPanel() {
	case PANEL_CONTRACTS:
		return NAVBAR_CONTRACTS
	case PANEL_SHIPS:
		return NAVBAR_SHIPS
	default:
		return NAVBAR_MAIN
	}
}

func (a *App) Navbar() tview.Primitive {
	switch a.UIState.SelectedPanel() {
	case PANEL_CONTRACTS:
		return a.contractsList()
	case PANEL_SHIPS:
		return a.shipsList()
	default:
		return a.mainNavList()
	}
}

func (a *App) UpdateNavbar() {
	a.UIState.selectedNavbar = a.UIState.DeterminedNavbar()
	a.navbar = a.Navbar()
}

func (a *App) mainNavList() *tview.List {
	return tview.NewList().
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
			a.SetSelectedPanel(PANEL_SHIPS)
		}).
		AddItem("Quit", "", 'q', func() {
			api.Close()
			a.ui.Stop()
		})
}

func (a *App) contractsList() *tview.List {
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
	if option != '0' {
		menu.AddItem("Negotiate New Contract", "", 'n', func() {
			api.NegotiateContract(a.ShipState.SelectedShipSymbol())
			a.draw(true)
		})
	}
	return menu
}

func (a *App) shipsList() *tview.List {
	ships := state.Ships(false)

	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			a.SetSelectedPanel(PANEL_DASHBOARD)
			a.draw(true)
		})
	if a.ShipState.SelectedShipIndex() != 0 && len(ships) > 1 {
		menu.AddItem("Load Previous Ship", "", 'p', func() {
			index := a.ShipState.SelectedShipIndex() - 1
			app.ShipState.SelectShip(index, ships[index].Symbol)
			a.draw(true)
		})
	}
	if a.ShipState.SelectedShipIndex() < len(ships)-1 {
		menu.AddItem("Load Next Ship", "", 'n', func() {
			index := a.ShipState.SelectedShipIndex() + 1
			app.ShipState.SelectShip(index, ships[index].Symbol)
			a.draw(true)
		})
	}

	return menu
}
