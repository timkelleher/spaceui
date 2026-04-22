package ui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/state"
)

const (
	MENU_SHIP_DETAIL = "ship_detail"
	PAGE_SHIP_DETAIL = "ship_detail"
)

type ShipDetailMenu struct {
}

func (sdm ShipDetailMenu) ID() string {
	return MENU_SHIP_DETAIL
}

func (sdm ShipDetailMenu) Menu() *tview.List {
	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			state.SetActivePage(PAGE_SHIPS_LIST)
			state.DeselectShip()
			app.draw(true)
		})

	selectedShip := state.SelectedShip()
	if state.HasActiveShip() && state.ActiveShip().Symbol == selectedShip.Symbol {
		switch selectedShip.Nav.Status {
		case "DOCKED":
			menu.AddItem("Orbit Ship", "", 'o', func() {
				api.OrbitShip(selectedShip.Symbol)

				state.Queue(state.DATA_SHIPS)
				state.SetActivePage(PAGE_SHIPS_LIST)
				app.draw(true)
			})
		case "IN_ORBIT":
			menu.AddItem("Dock Ship", "", 'd', func() {
				api.DockShip(selectedShip.Symbol)

				state.Queue(state.DATA_SHIPS)
				state.SetActivePage(PAGE_SHIPS_LIST)
				app.draw(true)
			})
		}

		menu.AddItem("Deactivate Ship", "", 'd', func() {
			state.DeactivateShip()
			app.draw(true)
		})
	} else {
		menu.AddItem("Activate Ship", "", 'a', func() {
			state.SetActiveShip(selectedShip)
			app.draw(true)
		})
	}

	return menu
}

type ShipDetailPage struct {
}

func (sdp ShipDetailPage) ID() string {
	return PAGE_SHIP_DETAIL
}

func (sdp ShipDetailPage) Menu() Menu {
	return ShipDetailMenu{}
}

func (sdp ShipDetailPage) RequiredData() []string {
	return []string{state.DATA_SHIPS}
}

// TODO: don't need this anymore since moving away from indexes?
func selectedShipIndex(selectedShip *api.Ship, ships []*api.Ship) int {
	for i, ship := range ships {
		if selectedShip.Symbol == ship.Symbol {
			return i
		}
	}
	return -1
}

func (sdp ShipDetailPage) Content() string {
	ships := state.Ships(false)

	if len(ships) == 0 {
		return "[red]Error: no available ships![-]"
	}

	selectedShip := state.SelectedShip()
	content := fmt.Sprintf("[red]Ship %s[-]\n", selectedShip.Symbol)
	content += fmt.Sprintf("[orange]Frame:[-] %s | [orange]Reactor:[-] %s | [orange]Engine:[-] %s\n", selectedShip.Frame.Name, selectedShip.Reactor.Name, selectedShip.Engine.Name)
	content += "----- Registration -----\n"
	content += fmt.Sprintf("[blue]Name:[-] %s | [blue]Faction Symbol:[-] %s | [blue]Role:[-] %s\n", selectedShip.Registration.Name, selectedShip.Registration.FactionSymbol, selectedShip.Registration.Role)
	content += "----- Nav -----\n"
	content += fmt.Sprintf("[yellow]System:[-] %s | [yellow]Waypoint:[-] %s\n", selectedShip.Nav.SystemSymbol, selectedShip.Nav.WaypointSymbol)
	content += fmt.Sprintf("[yellow]Status:[-] %s | [yellow]Flight Mode:[-] %s\n", selectedShip.Nav.Status, selectedShip.Nav.FlightMode)
	content += "----- Cargo -----\n"
	content += fmt.Sprintf("[green]Capacity:[-] %d | [green]Units:[-] %d\n", selectedShip.Cargo.Capacity, selectedShip.Cargo.Units)
	content += "----- Fuel -----\n"
	content += fmt.Sprintf("[green]Capacity:[-] %d | [green]Current:[-] %d | [green]Consumed:[-] %d\n", selectedShip.Fuel.Capacity, selectedShip.Fuel.Current, selectedShip.Fuel.Consumed.Amount)
	return content
}
