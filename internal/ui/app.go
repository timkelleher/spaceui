package ui

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/events"
	"github.com/timkelleher/spaceui/internal/logger"
	"github.com/timkelleher/spaceui/internal/state"
)

// /////////////////////////////////////
// App
// /////////////////////////////////////
var app *App

type App struct {
	mu sync.Mutex

	ui   *tview.Application
	grid *tview.Grid

	pages    map[string]Page
	homePage Page

	currentMenuID string
	uiMenu        *tview.List
	uiPage        *tview.TextView
	uiStatusBar   *tview.TextView
	uiFooter      *tview.TextView
}

func (a *App) registerPages() {
	a.pages = make(map[string]Page)

	a.registerPage(DashboardPage{}, true)

	a.registerPage(ContractsPage{}, false)

	a.registerPage(ShipsListPage{}, false)
	a.registerPage(ShipDetailPage{}, false)

	a.registerPage(WaypointsListPage{}, false)
	a.registerPage(WaypointDetailPage{}, false)
	a.registerPage(WaypointMarketplacePage{}, false)
	a.registerPage(WaypointShipyardPage{}, false)
}

func (a *App) registerPage(p Page, isHome bool) {
	a.pages[p.ID()] = p
	if isHome {
		a.homePage = p
	}
}

func (a *App) getPage(id string) Page {
	for _, page := range a.pages {
		if page.ID() == id {
			return page
		}
	}
	return nil
}

func NewApp() *App {
	if app != nil {
		return app
	}

	api.SetApiKey(os.Getenv("SPACE_TRADERS_API_KEY"))

	app = &App{
		mu: sync.Mutex{},

		ui: tview.NewApplication(),
		grid: tview.NewGrid().
			SetRows(0, 1, 10).
			SetColumns(35, 0).
			SetBorders(true),
		uiPage: tview.NewTextView().
			SetDynamicColors(true).
			SetWrap(true).
			SetScrollable(true),
		uiStatusBar: tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignCenter),
		uiFooter: tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignLeft),
	}

	app.registerPages()
	state.SetActivePage(app.homePage.ID())
	app.uiPage.
		SetText(app.homePage.Content()).
		SetChangedFunc(func() {
			app.ui.Draw()
		})

	return app
}

func (a *App) Run() {
	// TODO: move this entirely to state?
	loc, _ := time.LoadLocation(os.Getenv("TIMEZONE"))
	state.SetLoc(loc)

	// Redraw every second
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

func (a *App) getCurrentMenuID() string {
	return a.currentMenuID
}

func (a *App) setCurrentMenuID(id string) {
	a.currentMenuID = id
}

func (a *App) draw(updateMenu bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	ev := events.GetEvents()
	for _, event := range ev {
		if updateMenu {
			break
		}
		if a.shouldRedrawMenuFromEvent(event) {
			updateMenu = true
		}
	}

	currentPage := a.getPage(state.ActivePage())
	if currentPage == nil {
		currentPage = a.getPage(PAGE_DASHBOARD)
		state.SetActivePage(currentPage.ID())
	}

	// Menu will be updated, so we should re-draw it
	if currentPage.Menu().ID() != a.getCurrentMenuID() {
		updateMenu = true
	}

	// Redraw the menu if forced or if the menu has changed since last draw
	if updateMenu {
		a.grid.RemoveItem(a.uiMenu)
		a.uiMenu = currentPage.Menu().Menu()
		a.setCurrentMenuID(currentPage.Menu().ID())
		a.grid.AddItem(a.uiMenu, 0, 0, 1, 1, 0, 0, false)
		a.ui.SetFocus(a.uiMenu)
	}

	a.grid.RemoveItem(a.uiPage)
	if state.AreRefreshing(currentPage.RequiredData()) {
		for _, datatype := range currentPage.RequiredData() {
			if !state.Fresh(datatype) {
				state.Queue(datatype)
			}
		}
		a.uiPage.SetText(loadingContent(currentPage.RequiredData()))
	} else {
		a.uiPage.SetText(currentPage.Content())
	}

	a.grid.AddItem(a.uiPage, 0, 1, 1, 1, 0, 0, false)

	a.grid.RemoveItem(a.uiStatusBar)
	a.uiStatusBar.SetText(statusBarContent())
	a.grid.AddItem(a.uiStatusBar, 1, 0, 1, 2, 0, 0, false)

	a.uiFooter.SetText(logger.Logs())
	a.grid.RemoveItem(a.uiFooter)
	a.grid.AddItem(a.uiFooter, 2, 0, 1, 2, 0, 0, false)
}

func (a *App) shouldRedrawMenuFromEvent(event string) bool {
	switch event {
	case events.EVENT_LOAD_SHIPS_COMPLETE:
		return a.currentMenuID == MENU_SHIPS_LIST || a.currentMenuID == MENU_SHIP_DETAIL
	case events.EVENT_LOAD_WAYPOINTS_COMPLETE:
		return a.currentMenuID == MENU_WAYPOINTS_LIST || a.currentMenuID == MENU_SHIP_DETAIL
	}
	return false
}

func loadingContent(datatypes []string) string {
	count := 0
	for _, datatype := range datatypes {
		count += state.RefreshStatus(datatype)
	}

	objNames := strings.Join(datatypes, ", ")
	if count == 0 {
		return fmt.Sprintf("[grey]Loading %s...[-]\n", objNames)
	}
	return fmt.Sprintf("[grey]Loading %d %s...[-]\n", count, objNames)
}

// /////////////////////////////////////
// Main Menu
// /////////////////////////////////////
const MENU_MAIN = "main"

type MainMenu struct {
}

func (mm MainMenu) ID() string {
	return MENU_MAIN
}

func (mm MainMenu) Menu() *tview.List {
	menu := tview.NewList().
		AddItem("Dashboard", "", 'd', func() {
			state.SetActivePage(PAGE_DASHBOARD)
		}).
		AddItem("Contracts", "", 'c', func() {
			state.SetActivePage(PAGE_CONTRACTS)
		}).
		AddItem("Ships", "", 's', func() {
			state.SetActivePage(PAGE_SHIPS_LIST)
		})
		//AddItem("Systems", "", 'y', func() {
		//	state.SetActivePage(PAGE_SYSTEMS)
		//})

	if state.HasActiveShip() {
		menu.AddItem("Waypoints", "", 'w', func() {
			state.SetActivePage(PAGE_WAYPOINTS_LIST)
		})
	}

	menu.AddItem("Quit", "", 'q', func() {
		api.Close()
		app.ui.Stop()
	})

	return menu
}
