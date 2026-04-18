package ui

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
	"github.com/timkelleher/spaceui/internal/logger"
	"github.com/timkelleher/spaceui/internal/state"
)

type layout struct {
	app  *tview.Application
	grid *tview.Grid

	currentTab     string
	currentNavbar  string
	currentControl string

	dataPaneData string

	navbar      *tview.List
	dataPane    *tview.TextView
	controlPane tview.Primitive
	footerPane  *tview.TextView
}

func (l *layout) isCurrentTabMainMenu() bool {
	switch l.currentTab {
	case "dashboard", "agent", "contracts":
		return true
	default:
		return false
	}
}

func (l *layout) setCurrentTab(id string) {
	l.currentTab = id
}

///////////////////
// Navbar
///////////////////

func (l *layout) determineNavbar() string {
	if l.isCurrentTabMainMenu() {
		return "main"
	}

	return "ships"
}

func (l *layout) navbarList() *tview.List {
	if l.isCurrentTabMainMenu() {
		return l.getMainNavList()
	}

	return l.getShipsSubList()
}

func (l *layout) getMainNavList() *tview.List {
	return tview.NewList().
		AddItem("Dashboard", "", 'd', func() {
			l.currentTab = "dashboard"
			l.setDataPaneData("")
			logger.Info("setting current tab to dashboard")
		}).
		AddItem("Agent", "", 'a', AgentAction).
		AddItem("Contracts", "", 'c', ContractsAction).
		AddItem("Ships", "", 's', ShipsAction).
		AddItem("Quit", "", 'q', func() {
			l.app.Stop()
		})
}

func (l *layout) getShipsSubList() *tview.List {
	return tview.NewList().
		AddItem("Back", "", 'b', func() {
			l.currentTab = "dashboard"
			l.setDataPaneData("")
			logger.Debug("setting current tab to dashboard")

			//grid.AddItem(sidebarList, 0, 0, 1, 1, 0, 0, true)
			//grid.RemoveItem(submenu)
			//app.SetFocus(sidebarList)
			//logger.Debug("setting current menu to main menu")

		})
}

func (l *layout) isNavbarFocused() bool {
	if l.isCurrentTabMainMenu() {
		return true
	}
	return true
}

///////////////////
// Control
///////////////////

func (l *layout) determineControl() string {
	if l.isCurrentTabMainMenu() {
		return "main"
	}

	return "ships"
}

func (l *layout) getControlPane() tview.Primitive {
	if l.isCurrentTabMainMenu() {
		return getMainControlPane()
	}

	return l.getShipsControlList()
}

func getMainControlPane() *tview.TextView {
	now := fmt.Sprintf("[green]%s[-]", time.Now().Format(time.RFC1123))
	agentSymbol := fmt.Sprintf("[red]%s[-]", state.Get("agent.symbol"))
	credits := fmt.Sprintf("[blue]%s[-]", state.Get("agent.credits"))

	return tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText(now + "\t" + agentSymbol + "\t" + credits)
}

func (l *layout) getShipsControlList() *tview.List {
	return tview.NewList().
		AddItem("Ships!!", "", 'd', func() {
			l.currentTab = "dashboard"
			logger.Info("setting current tab to dashboard")
		}).
		AddItem("Another!!", "", 'a', func() {
			l.currentTab = "dashboard"
			logger.Info("setting current tab to dashboard")
		}).
		SetOffset(0, 0)
}

func (l *layout) isControlFocused() bool {
	if l.isCurrentTabMainMenu() {
		return false
	}
	return false
}

// /////////////////
// Data
// /////////////////
func (l *layout) setDataPaneData(data string) {
	l.dataPane.SetText(data)
}

func NewLayout() layout {
	layout := layout{
		app: tview.NewApplication(),
		grid: tview.NewGrid().
			SetRows(0, 1, 10).
			SetColumns(30, 0).
			SetBorders(true),
		currentTab: "dashboard",

		controlPane: tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignCenter).
			SetText("Loading..."),
		footerPane: tview.NewTextView().
			SetDynamicColors(true).
			SetTextAlign(tview.AlignLeft).
			SetText(""),
	}
	layout.navbar = layout.getMainNavList()
	layout.dataPane = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetChangedFunc(func() {
			layout.app.Draw()
		})
	return layout
}

func (l *layout) start() {
	go func() {
		for range time.Tick(time.Second) {
			l.app.QueueUpdateDraw(func() {
				l.draw()
			})
		}
	}()

	if err := l.app.SetRoot(l.grid, true).Run(); err != nil {
		panic(err)
	}
}

func (l *layout) draw() {
	if l.currentNavbar != l.determineNavbar() {
		logger.Info(fmt.Sprintf("redrawing navbar from %s to %s", l.currentNavbar, l.determineNavbar()))
		l.currentNavbar = l.determineNavbar()
		l.grid.RemoveItem(l.navbar)
		l.navbar = l.navbarList()
		l.grid.AddItem(l.navbar, 0, 0, 1, 1, 0, 0, false)
		if l.isNavbarFocused() {
			l.app.SetFocus(l.navbar)
		}
	}

	l.grid.RemoveItem(l.dataPane)
	l.grid.AddItem(l.dataPane, 0, 1, 1, 1, 0, 0, false)

	if l.currentControl != l.determineControl() {
		l.currentControl = l.determineControl()
		l.grid.RemoveItem(l.controlPane)
		l.controlPane = l.getControlPane()
		l.grid.AddItem(l.controlPane, 1, 0, 1, 2, 0, 0, false)
		if l.isControlFocused() {
			l.app.SetFocus(l.controlPane)
		}
	}

	l.footerPane.SetText(logger.Logs())
	l.grid.RemoveItem(l.footerPane)
	l.grid.AddItem(l.footerPane, 2, 0, 1, 2, 0, 0, false)
}
