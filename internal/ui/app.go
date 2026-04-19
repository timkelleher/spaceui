package ui

import (
	"time"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/logger"
)

var app App

type App struct {
	ShipState
	UIState

	ui   *tview.Application
	grid *tview.Grid

	currentNavbar string

	navbar     tview.Primitive
	dataPane   *tview.TextView
	statusPane tview.Primitive
	footerPane *tview.TextView
}

func (a *App) PanelData() string {
	return panels[a.UIState.selectedPanel]()
}

func NewApp() App {
	app := App{
		UIState: UIState{selectedPanel: PANEL_DASHBOARD},

		ui: tview.NewApplication(),
		grid: tview.NewGrid().
			SetRows(0, 1, 10).
			SetColumns(30, 0).
			SetBorders(true),

		statusPane: tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignCenter).
			SetText("Loading..."),
		footerPane: tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignLeft).
			SetText(""),
	}
	app.dataPane = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetChangedFunc(func() {
			app.ui.Draw()
		})
	return app
}

func (a *App) Run() {
	go func() {
		for range time.Tick(time.Second) {
			a.ui.QueueUpdateDraw(func() {
				a.draw(false)
			})
		}
	}()

	if err := a.ui.SetRoot(a.grid, true).Run(); err != nil {
		panic(err)
	}
}

func (a *App) draw(forceNavUpdate bool) {
	if forceNavUpdate || a.UIState.NavbarChanged() {
		a.grid.RemoveItem(a.navbar)
		a.UpdateNavbar()
		a.grid.AddItem(a.navbar, 0, 0, 1, 1, 0, 0, false)
		a.ui.SetFocus(a.navbar)
	}

	a.grid.RemoveItem(a.dataPane)
	a.dataPane.SetText(a.PanelData())
	a.grid.AddItem(a.dataPane, 0, 1, 1, 1, 0, 0, false)

	a.grid.RemoveItem(a.statusPane)
	a.statusPane = a.statusBarContent()
	a.grid.AddItem(a.statusPane, 1, 0, 1, 2, 0, 0, false)

	a.footerPane.SetText(logger.Logs())
	a.grid.RemoveItem(a.footerPane)
	a.grid.AddItem(a.footerPane, 2, 0, 1, 2, 0, 0, false)

	// select current ship if we don't have one yet
	if !a.ShipState.HasSelectedShip() {
		a.ShipState.SelectFromCache()
	}
}
