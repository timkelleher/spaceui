package ui

import (
	"fmt"
	"time"

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

	if !state.HasSelectedShip() {
		return menu
	}

	selectedShip := state.SelectedShip()
	if state.IsSelectedActiveShip() {
		switch selectedShip.Nav.Status {
		case "DOCKED":
			menu.AddItem("Orbit Ship", "", 'o', func() {
				api.OrbitShip(selectedShip.Symbol)

				// Delay data refresh
				go func() {
					state.MarkStale(state.DATA_SHIPS)

					time.Sleep(2 * time.Second)
					state.Queue(state.DATA_SHIPS)
				}()

				state.SetActivePage(PAGE_SHIPS_LIST)
				app.draw(true)
			})

			if selectedShip.Fuel.Current != selectedShip.Fuel.Capacity {
				menu.AddItem("Refuel Ship", "", 'r', func() {
					api.RefuelShip(selectedShip.Symbol)

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
		case "IN_ORBIT":
			menu.AddItem("Dock Ship", "", 'd', func() {
				api.DockShip(selectedShip.Symbol)

				// Delay data refresh
				go func() {
					time.Sleep(2 * time.Second)
					state.Queue(state.DATA_SHIPS)
				}()

				app.draw(true)
			})

			if selectedShip.CanMine() && !selectedShip.OnCooldown() {
				menu.AddItem("Extract", selectedShip.Nav.WaypointSymbol, 'e', func() {
					api.ExtractShip(selectedShip.Symbol)

					// Delay data refresh
					go func() {
						state.MarkStale(state.DATA_SHIPS)

						time.Sleep(2 * time.Second)
						state.Queue(state.DATA_SHIPS)
					}()

					app.draw(true)
				})
			}
		}

		if !state.AnyRefreshing() {
			menu.AddItem("Deactivate Ship", "", 'd', func() {
				if !state.AnyRefreshing() {
					state.DeactivateShip()

					state.SetActivePage(PAGE_SHIPS_LIST)
					app.draw(true)
				}
			})
		}
	} else if !state.AnyRefreshing() {
		menu.AddItem("Activate Ship", "", 'a', func() {
			if !state.AnyRefreshing() {
				state.SetActiveShip(selectedShip)

				state.SetActivePage(PAGE_SHIPS_LIST)
				app.draw(true)
			}
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

func (sdp ShipDetailPage) Content() string {
	ships := state.Ships(false)

	if len(ships) == 0 {
		return "[red]Error: no available ships![-]"
	}
	if !state.HasSelectedShip() {
		return "[red]Error: invalid selected ship state![-]"
	}

	selectedShip := state.SelectedShip()
	content := fmt.Sprintf("[red]Ship %s[-]\n", selectedShip.Symbol)
	content += fmt.Sprintf("[orange]Frame:[-] %s | [orange]Reactor:[-] %s | [orange]Engine:[-] %s\n", selectedShip.Frame.Name, selectedShip.Reactor.Name, selectedShip.Engine.Name)
	content += "----- Registration -----\n"
	content += fmt.Sprintf("[blue]Name:[-] %s | [blue]Faction Symbol:[-] %s | [blue]Role:[-] %s\n", selectedShip.Registration.Name, selectedShip.Registration.FactionSymbol, selectedShip.Registration.Role)
	content += "----- Nav -----\n"
	content += fmt.Sprintf("[yellow]System:[-] %s | [yellow]Waypoint:[-] %s\n", selectedShip.Nav.SystemSymbol, selectedShip.Nav.WaypointSymbol)
	content += fmt.Sprintf("[yellow]Status:[-] %s | [yellow]Flight Mode:[-] %s\n", selectedShip.Nav.Status, selectedShip.Nav.FlightMode)

	if len(selectedShip.Modules) > 0 {
		content += "----- Modules -----\n"
		for _, module := range selectedShip.Modules {
			content += fmt.Sprintf("[yellow]Name:[-] %s | [yellow]Capacity:[-] %d\n", module.Name, module.Capacity)
			for _, requirement := range module.Requirements {
				content += fmt.Sprintf("\t%d Crew | %d Power | %d Slots\n", requirement.Crew, requirement.Power, requirement.Slots)
			}
		}
	}
	if len(selectedShip.Mounts) > 0 {
		content += "----- Mounts -----\n"
		for _, mount := range selectedShip.Mounts {
			content += fmt.Sprintf("[yellow]Name:[-] %s | [yellow]Strength:[-] %d\n", mount.Name, mount.Strength)
			for _, requirement := range mount.Requirements {
				content += fmt.Sprintf("\t%d Crew | %d Power\n", requirement.Crew, requirement.Power)
			}
		}
	}

	content += "----- Cargo -----\n"
	content += fmt.Sprintf("[green]Capacity:[-] %d | [green]Units:[-] %d\n", selectedShip.Cargo.Capacity, selectedShip.Cargo.Units)
	for _, item := range selectedShip.Cargo.Inventory {
		content += fmt.Sprintf("\t%d [green]%s[-]\n", item.Units, item.Name)
	}
	content += "----- Fuel -----\n"
	content += fmt.Sprintf("[green]Capacity:[-] %d | [green]Current:[-] %d | [green]Consumed:[-] %d\n", selectedShip.Fuel.Capacity, selectedShip.Fuel.Current, selectedShip.Fuel.Consumed.Amount)
	content += "----- Cooldown -----\n"
	content += fmt.Sprintf("[purple]Remaining Cooldown[-] %d/%d sec\n", selectedShip.Cooldown.RemainingSeconds, selectedShip.Cooldown.TotalSeconds)
	return content
}
