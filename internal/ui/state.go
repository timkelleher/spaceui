package ui

import (
	"fmt"

	"github.com/timkelleher/spaceui/internal/api"
	"github.com/timkelleher/spaceui/internal/logger"
	"github.com/timkelleher/spaceui/internal/state"
)

const (
	PANEL_DASHBOARD                = "dashboard"
	PANEL_AGENT                    = "agent"
	PANEL_CONTRACTS                = "contracts"
	PANEL_SHIPS_LIST               = "ships_list"
	PANEL_SHIP_DETAIL              = "ship_detail"
	PANEL_SYSTEMS                  = "systems"
	PANEL_WAYPOINTS_LIST           = "waypoints_list"
	PANEL_WAYPOINT_DETAIL          = "waypoint_detail"
	PANEL_WAYPOINT_AVAILABLE_SHIPS = "waypoint_available_ships"

	MENU_MAIN                     = "main"
	MENU_CONTRACTS                = "contracts"
	MENU_SHIPS_LIST               = "ships_list"
	MENU_SHIP_DETAIL              = "ship_detail"
	MENU_SYSTEMS                  = "systems"
	MENU_WAYPOINTS_LIST           = "waypoints_list"
	MENU_WAYPOINT_DETAIL          = "waypoint_detail"
	MENU_WAYPOINT_AVAILABLE_SHIPS = "waypoint_available_ships"
)

/////////////////////////////
// UIState
/////////////////////////////

type UIState struct {
	selectedMenu  string
	selectedPanel string
}

func (us *UIState) SelectedPanel() string {
	return us.selectedPanel
}

// TODO: safety or fallback?
func (us *UIState) SetSelectedPanel(id string) {
	us.selectedPanel = id
}

/////////////////////////////
// GameState
/////////////////////////////

type GameState struct {
	selectedShipIndex int
	activeShipIndex   int
	activeShipSymbol  string

	selectedWaypointSymbol string
	waypointFilterType     string
	waypointFilterName     string
}

func NewGameState() GameState {
	return GameState{
		selectedShipIndex: -1,
		activeShipIndex:   -1,
	}
}

/////////////////////////////
// Ships
/////////////////////////////

func (gs *GameState) SelectShip(index int) {
	gs.selectedShipIndex = index
}

func (gs *GameState) DeselectShip() {
	gs.selectedShipIndex = -1
}

func (gs *GameState) SelectedShipIndex() int {
	return gs.selectedShipIndex
}

func (gs *GameState) HasActiveShip() bool {
	return gs.activeShipSymbol != ""
}

func (gs *GameState) ActiveShipSymbol() string {
	return gs.activeShipSymbol
}

func (gs *GameState) ActiveShip() *api.Ship {
	if !gs.HasActiveShip() {
		return nil
	}
	ships := state.Ships(false)
	for _, ship := range ships {
		if ship.Symbol == gs.ActiveShipSymbol() {
			return &ship
		}
	}
	return nil
}

func (gs *GameState) ActivateShip(index int, symbol string) {
	gs.selectedShipIndex = index
	gs.activeShipSymbol = symbol

	state.Reset("waypoints")
	logger.Info(fmt.Sprintf("STATE active ship [green]%s[-] (%d)", symbol, index))
}

func (gs *GameState) ActivateFromCache() {
	ships := state.Ships(false)
	if len(ships) > 0 {
		gs.ActivateShip(0, ships[0].Symbol)
	}
}

func (gs *GameState) DeactivateShip() {
	gs.selectedShipIndex = -1
	gs.activeShipSymbol = ""
}

/////////////////////////////
// Waypoints
/////////////////////////////

func (gs *GameState) ActiveWaypoint() *api.Waypoint {
	ship := gs.ActiveShip()
	for _, waypoint := range state.Waypoints(ship.Nav.SystemSymbol) {
		if waypoint.Symbol == gs.selectedWaypointSymbol {
			return &waypoint
		}
	}
	return nil
}

func (gs *GameState) ApplyWaypointFilter(waypointType, name string) {
	gs.waypointFilterType = waypointType
	gs.waypointFilterName = name
}
