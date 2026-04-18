package ui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/logger"
	"github.com/timkelleher/spaceui/internal/state"
)

var (
	currentShipIndex int
	currentShipName  string
)

func ShipsAction() {
	app.setCurrentTab("ships")
	app.setDataPaneData("Loading...")
	logger.Info("setting current tab to ships")

	go func() {
		//ships := state.Ships()
		app.setDataPaneData(ShipsTabContent())

		subMenu := tview.NewList()
		subMenu.AddItem("Back", "", 'b', func() {
			app.setCurrentTab("dashboard")
			logger.Debug("setting current tab to dashboard")

			//grid.AddItem(mainMenu, 0, 0, 1, 1, 0, 0, true)
			//grid.RemoveItem(submenu)
			//app.SetFocus(mainMenu)
			//currentMenu = mainMenu
			//logger.Debug("setting current menu to main menu")
		}).
			AddItem("Current Ship", currentShipName, ' ', nil)

		/*
			if currentShipIndex != 0 && len(ships) > 1 {
				subMenu.AddItem("Load Previous Ship", "", 'p', func() {
					currentTab = "ships"
					currentShipIndex--
					main.SetText(ShipsTabContent())

					grid.AddItem(subMenu, 0, 0, 1, 1, 0, 0, true)
					grid.RemoveItem(currentMenu)
					app.SetFocus(subMenu)
					currentMenu = subMenu
					logger.Debug("setting current menu to ship submenu")
				})
			}
			if currentShipIndex < len(ships)-1 {
				subMenu.AddItem("Load Next Ship", "", 'p', func() {
					currentTab = "ships"
					currentShipIndex++
					main.SetText(ShipsTabContent())

					grid.AddItem(subMenu, 0, 0, 1, 1, 0, 0, true)
					grid.RemoveItem(currentMenu)
					app.SetFocus(subMenu)
					currentMenu = subMenu
					logger.Debug("setting current menu to ship submenu")
				})
			}

			grid.RemoveItem(currentMenu)
			grid.AddItem(subMenu, 0, 0, 1, 1, 0, 0, true)
			app.SetFocus(subMenu)
			currentMenu = subMenu
			logger.Debug("setting current menu to ship submenu")
		*/
	}()
}

func ShipsTabContent() string {
	ships := state.Ships()

	content := fmt.Sprintf("[yellow]Number of ships owned:[-] %d\n", len(ships))
	if len(ships) == 0 {
		return content
	}

	ship := ships[currentShipIndex]
	currentShipName = ship.Symbol

	content += fmt.Sprintf("Current Ship selected: [yellow]%s[-]\n", ships[currentShipIndex].Symbol)
	content += "----------\n"
	content += fmt.Sprintf("[red]Ship %s[-]\n", ship.Symbol)
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
	if currentShipIndex != len(ships)-1 {
		content += "\n"
	}

	return content
}
