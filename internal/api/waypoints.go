package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	GET_WAYPOINTS_ENDPOINT                = "/systems/%s/waypoints"
	GET_WAYPOINT_AVAILABLE_SHIPS_ENDPOINT = "/systems/%s/waypoints/%s/shipyard"
)

type WaypointsResponse struct {
	Waypoints []Waypoint       `json:"data"`
	Meta      ListMetaResponse `json:"meta"`
}

type Waypoint struct {
	Symbol       string    `json:"symbol"`
	Type         string    `json:"type"`
	SystemSymbol string    `json:"systemSymbol"`
	X            int       `json:"x"`
	Y            int       `json:"y"`
	Orbitals     []Orbital `json:"orbitals"`
	Traits       []struct {
		Symbol      string `json:"symbol"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"traits"`
	IsUnderConstruction bool    `json:"isUnderConstruction"`
	Faction             Faction `json:"faction"`
	Modifiers           []any   `json:"modifiers"`
	Chart               struct {
		WaypointSymbol string    `json:"waypointSymbol"`
		SubmittedBy    string    `json:"submittedBy"`
		SubmittedOn    time.Time `json:"submittedOn"`
	} `json:"chart"`
}

func (w Waypoint) IsShipyard() bool {
	for _, trait := range w.Traits {
		if trait.Name == "Shipyard" {
			return true
		}
	}
	return false
}

func (w Waypoint) AllTraits() []string {
	var traits []string
	for _, trait := range w.Traits {
		traits = append(traits, trait.Name)
	}

	return traits
}

type Orbital struct {
	Symbol string `json:"symbol"`
}

func GetWaypoints(system string, page int) (*WaypointsResponse, ApiResult) {
	endpoint := fmt.Sprintf(GET_WAYPOINTS_ENDPOINT, system)

	var obj WaypointsResponse
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetQueryParams(map[string]string{"page": strconv.Itoa(page), "limit": "20"}).
		SetResult(&obj).
		SetError(&errResp).
		Get(url + endpoint)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodGet, endpoint+fmt.Sprintf(" (%d)", page), res)
	return &obj, res
}

type WaypointAvailableShipsResponse struct {
	Shipyard `json:"data"`
}

type Shipyard struct {
	Symbol    string `json:"symbol"`
	ShipTypes []struct {
		Type string `json:"type"`
	} `json:"shipTypes"`
	ModificationsFee int             `json:"modificationsFee"`
	Transactions     []interface{}   `json:"transactions"`
	Ships            []AvailableShip `json:"ships"`
}

type AvailableShip struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Supply        string `json:"supply"`
	PurchasePrice int    `json:"purchasePrice"`
	Frame         struct {
		Symbol         string `json:"symbol"`
		Name           string `json:"name"`
		Condition      int    `json:"condition"`
		Integrity      int    `json:"integrity"`
		Description    string `json:"description"`
		ModuleSlots    int    `json:"moduleSlots"`
		MountingPoints int    `json:"mountingPoints"`
		FuelCapacity   int    `json:"fuelCapacity"`
		Requirements   struct {
			Power int `json:"power"`
			Crew  int `json:"crew"`
		} `json:"requirements"`
		Quality int `json:"quality"`
	} `json:"frame"`
	Reactor struct {
		Symbol       string `json:"symbol"`
		Name         string `json:"name"`
		Condition    int    `json:"condition"`
		Integrity    int    `json:"integrity"`
		Description  string `json:"description"`
		PowerOutput  int    `json:"powerOutput"`
		Requirements struct {
			Crew int `json:"crew"`
		} `json:"requirements"`
		Quality int `json:"quality"`
	} `json:"reactor"`
	Engine struct {
		Symbol       string `json:"symbol"`
		Name         string `json:"name"`
		Condition    int    `json:"condition"`
		Integrity    int    `json:"integrity"`
		Description  string `json:"description"`
		Speed        int    `json:"speed"`
		Requirements struct {
			Power int `json:"power"`
			Crew  int `json:"crew"`
		} `json:"requirements"`
		Quality int `json:"quality"`
	} `json:"engine"`
	Modules []struct {
		Symbol       string `json:"symbol"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		Requirements struct {
			Power int `json:"power"`
			Crew  int `json:"crew"`
			Slots int `json:"slots"`
		} `json:"requirements"`
		Capacity int `json:"capacity,omitempty"`
	} `json:"modules"`
	Mounts []struct {
		Symbol       string `json:"symbol"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		Requirements struct {
			Power int `json:"power"`
			Crew  int `json:"crew"`
		} `json:"requirements"`
		Strength int `json:"strength"`
	} `json:"mounts"`
	Crew struct {
		Required int `json:"required"`
		Capacity int `json:"capacity"`
	} `json:"crew"`
	Activity string `json:"activity"`
}

func GetAvailableShips(system, waypoint string) (*WaypointAvailableShipsResponse, ApiResult) {
	endpoint := fmt.Sprintf(GET_WAYPOINT_AVAILABLE_SHIPS_ENDPOINT, system, waypoint)

	var obj WaypointAvailableShipsResponse
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&obj).
		SetError(&errResp).
		Get(url + endpoint)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodGet, endpoint, res)
	return &obj, res
}
