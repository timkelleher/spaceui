package ui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/state"
)

const (
	MENU_WAYPOINT_MARKETPLACE = "waypoint_marketplace"
	PAGE_WAYPOINT_MARKETPLACE = "waypoint_marketplace"
)

type WaypointMarketplaceMenu struct {
}

func (wmm WaypointMarketplaceMenu) ID() string {
	return MENU_WAYPOINT_MARKETPLACE
}

func (wmm WaypointMarketplaceMenu) Menu() *tview.List {
	menu := tview.NewList().
		AddItem("Back", "", 'b', func() {
			state.SetActivePage(PAGE_WAYPOINT_DETAIL)
			app.draw(true)
		})

	activeShip := state.ActiveShip()
	waypoint := state.SelectedWaypoint()

	if activeShip.Nav.WaypointSymbol == waypoint.Symbol {
		for _, item := range activeShip.Cargo.Inventory {
			menu.AddItem(fmt.Sprintf("Sell %s", item.Name), fmt.Sprintf("%d units", item.Units), 0, func() {
				api.SellCargo(activeShip.Symbol, item.Symbol, item.Units)

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
	}

	return menu
}

type WaypointMarketplacePage struct {
}

func (wmp WaypointMarketplacePage) ID() string {
	return PAGE_WAYPOINT_MARKETPLACE
}

func (wmp WaypointMarketplacePage) Menu() Menu {
	return WaypointMarketplaceMenu{}
}

func (wmp WaypointMarketplacePage) RequiredData() []string {
	return []string{state.DATA_WAYPOINTS}
}

func (wmp WaypointMarketplacePage) Content() string {
	waypoint := state.SelectedWaypoint()
	if waypoint == nil {
		return "[red]Error: no active waypoint![-]"
	}

	marketplace := state.Marketplace(waypoint)
	content := fmt.Sprintf("Marketplace on [blue]%s[-] with [yellow]%d[-] imports and [yellow]%d[-] exports\n", waypoint.Symbol, len(marketplace.Imports), len(marketplace.Exports))
	content += "----- Exchanges -----\n"
	for _, exchange := range marketplace.Exchange {
		content += fmt.Sprintf("[yellow]Type:[-]: %s | %s\n", exchange.Name, exchange.Description)
	}
	content += "----- Trade Goods -----\n"
	for _, good := range marketplace.TradeGoods {
		content += fmt.Sprintf("[orange]%s %s[-] [yellow]Supply:[-] %s [yellow]Volume:[-] %d [green]Purchase Price:[-] %d [green]Sell Price[-]: %d\n", good.Symbol, good.Type, good.Supply, good.TradeVolume, good.PurchasePrice, good.SellPrice)
	}
	content += "----- Transactions -----\n"
	for _, transaction := range marketplace.Transactions {
		content += fmt.Sprintf("[orange]%s[-] [yellow]Ship:[-] %s [orange]%s[-] [green]Units:[-] %d @ %d [green]Total:[-] %d\n", transaction.Type, transaction.ShipSymbol,
			transaction.TradeSymbol, transaction.Units, transaction.PricePerUnit, transaction.TotalPrice)
	}
	return content
}
