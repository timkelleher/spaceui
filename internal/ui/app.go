package ui

import (
	"os"
	"time"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/events"
	"github.com/timkelleher/spaceui/internal/logger"
	"github.com/timkelleher/spaceui/internal/state"
)

type App struct {
	GameState
	UIState

	ui   *tview.Application
	grid *tview.Grid

	loc *time.Location

	navbar        tview.Primitive
	dataPane      *tview.TextView
	statusBarPane *tview.TextView
	footerPane    *tview.TextView
}

func (a *App) PanelContent() string {
	switch a.UIState.SelectedPanel() {
	case PANEL_CONTRACTS:
		return a.contractsContent()
	case PANEL_SHIPS_LIST:
		return a.shipsListContent()
	case PANEL_SHIP_DETAIL:
		return a.shipDetailContent()
	case PANEL_SYSTEMS:
		return a.systemsContent()
	case PANEL_WAYPOINTS_LIST:
		return a.waypointsListContent()
	case PANEL_WAYPOINT_DETAIL:
		return a.waypointDetailContent()
	case PANEL_WAYPOINT_AVAILABLE_SHIPS:
		return a.waypointAvailableShipsContent()
	default:
		return a.dashboardContent()
	}
}

func (a *App) FormattedTime(t time.Time) string {
	loc := state.Loc()
	if loc == nil {
		return t.Format(time.RFC1123)
	}
	return t.In(loc).Format(time.RFC1123)
}

func NewApp() App {
	api.SetApiKey(os.Getenv("SPACE_TRADERS_API_KEY"))

	app := App{
		UIState:   UIState{selectedPanel: PANEL_DASHBOARD},
		GameState: NewGameState(),

		ui: tview.NewApplication(),
		grid: tview.NewGrid().
			SetRows(0, 1, 10).
			SetColumns(35, 0).
			SetBorders(true),
		dataPane: tview.NewTextView().
			SetDynamicColors(true).
			SetWrap(true).
			SetScrollable(true),
		statusBarPane: tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignCenter),
		footerPane: tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignLeft),
	}

	app.dataPane.SetChangedFunc(func() {
		app.ui.Draw()
	})
	return app
}

func (a *App) Run() {
	loc, _ := time.LoadLocation(os.Getenv("TIMEZONE"))
	state.SetLoc(loc)

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
	ev := events.GetEvents()
	for _, event := range ev {
		if event == events.EVENT_LOAD_WAYPOINTS_COMPLETE && a.UIState.selectedMenu == "waypoints_list" {
			forceNavUpdate = true
		}
	}

	if forceNavUpdate || a.UIState.MenuChanged() {
		a.grid.RemoveItem(a.navbar)
		a.UpdateMenu()
		a.grid.AddItem(a.navbar, 0, 0, 1, 1, 0, 0, false)
		a.ui.SetFocus(a.navbar)
	}

	a.grid.RemoveItem(a.dataPane)
	a.dataPane.SetText(a.PanelContent())
	a.grid.AddItem(a.dataPane, 0, 1, 1, 1, 0, 0, false)

	a.grid.RemoveItem(a.statusBarPane)
	a.statusBarPane.SetText(a.statusBarContent())
	a.grid.AddItem(a.statusBarPane, 1, 0, 1, 2, 0, 0, false)

	a.footerPane.SetText(logger.Logs())
	a.grid.RemoveItem(a.footerPane)
	a.grid.AddItem(a.footerPane, 2, 0, 1, 2, 0, 0, false)
}
