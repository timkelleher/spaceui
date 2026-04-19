package ui

import (
	"fmt"

	"github.com/timkelleher/spaceui/internal/logger"
	"github.com/timkelleher/spaceui/internal/state"
)

const (
	PANEL_DASHBOARD = "dashboard"
	PANEL_AGENT     = "agent"
	PANEL_CONTRACTS = "contracts"
	PANEL_SHIPS     = "ships"
	PANEL_SYSTEMS   = "systems"

	NAVBAR_MAIN      = "main"
	NAVBAR_CONTRACTS = "contracts"
	NAVBAR_SHIPS     = "ships"
	NAVBAR_SYSTEMS   = "systems"
)

type UIState struct {
	selectedNavbar string
	selectedPanel  string
}

func (us *UIState) SelectedPanel() string {
	return us.selectedPanel
}

// TODO: safety or fallback?
func (us *UIState) SetSelectedPanel(id string) {
	us.selectedPanel = id
}

type ShipState struct {
	selectedShipIndex  int
	selectedShipSymbol string
}

func (ss *ShipState) SelectedShipIndex() int {
	return ss.selectedShipIndex
}

func (ss *ShipState) SelectedShipSymbol() string {
	return ss.selectedShipSymbol
}

func (ss *ShipState) SelectShip(index int, symbol string) {
	ss.selectedShipIndex = index
	ss.selectedShipSymbol = symbol

	logger.Info(fmt.Sprintf("STATE selected ship [green]%s[-] (%d)", symbol, index))
}

func (ss *ShipState) HasSelectedShip() bool {
	return ss.selectedShipSymbol != ""
}

func (ss *ShipState) SelectFromCache() {
	ships := state.Ships(false)
	if len(ships) > 0 {
		ss.SelectShip(0, ships[0].Symbol)
	}
}
