package api

import (
	"fmt"
	"net/http"
	"time"
)

const (
	GET_SHIPS_ENDPOINT  = "/my/ships"
	ORBIT_SHIP_ENDPOINT = "/my/ships/%s/orbit"
	DOCK_SHIP_ENDPOINT  = "/my/ships/%s/dock"
)

type ShipsResponse struct {
	Ships    []Ship   `json:"data,omitempty"`
	ShipMeta ShipMeta `json:"meta,omitempty"`
}

type Ship struct {
	Symbol       string `json:"symbol"`
	Registration struct {
		Name          string `json:"name"`
		FactionSymbol string `json:"factionSymbol"`
		Role          string `json:"role"`
	} `json:"registration"`
	Nav struct {
		SystemSymbol   string `json:"systemSymbol"`
		WaypointSymbol string `json:"waypointSymbol"`
		Route          struct {
			Destination struct {
				Symbol       string `json:"symbol"`
				Type         string `json:"type"`
				SystemSymbol string `json:"systemSymbol"`
				X            int    `json:"x"`
				Y            int    `json:"y"`
			} `json:"destination"`
			Origin struct {
				Symbol       string `json:"symbol"`
				Type         string `json:"type"`
				SystemSymbol string `json:"systemSymbol"`
				X            int    `json:"x"`
				Y            int    `json:"y"`
			} `json:"origin"`
			DepartureTime time.Time `json:"departureTime"`
			Arrival       time.Time `json:"arrival"`
		} `json:"route"`
		Status     string `json:"status"`
		FlightMode string `json:"flightMode"`
	} `json:"nav"`
	Crew struct {
		Current  int    `json:"current"`
		Required int    `json:"required"`
		Capacity int    `json:"capacity"`
		Rotation string `json:"rotation"`
		Morale   int    `json:"morale"`
		Wages    int    `json:"wages"`
	} `json:"crew"`
	Frame struct {
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
		Strength int      `json:"strength"`
		Deposits []string `json:"deposits,omitempty"`
	} `json:"mounts"`
	Cargo struct {
		Capacity  int             `json:"capacity"`
		Units     int             `json:"units"`
		Inventory []ShipInventory `json:"inventory"`
	} `json:"cargo"`
	Fuel struct {
		Current  int `json:"current"`
		Capacity int `json:"capacity"`
		Consumed struct {
			Amount    int       `json:"amount"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"consumed"`
	} `json:"fuel"`
	Cooldown struct {
		ShipSymbol       string `json:"shipSymbol"`
		TotalSeconds     int    `json:"totalSeconds"`
		RemainingSeconds int    `json:"remainingSeconds"`
	} `json:"cooldown"`
}

type ShipInventory struct {
	Amount int `json:"amount"`
}

type ShipMeta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func GetShips() (*ShipsResponse, ApiResult) {
	var obj ShipsResponse
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&obj).
		SetError(&errResp).
		Get(url + GET_SHIPS_ENDPOINT)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodGet, GET_SHIPS_ENDPOINT, res)
	return &obj, res
}

func BuyShip(shipType, waypointSymbol string) ApiResult {
	payloadStruct := struct {
		ShipType       string `json:"shipType"`
		WaypointSymbol string `json:"waypointSymbol"`
	}{
		ShipType:       shipType,
		WaypointSymbol: waypointSymbol,
	}

	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetError(&errResp).
		SetBody(payloadStruct).
		Post(url + GET_SHIPS_ENDPOINT)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodPost, GET_SHIPS_ENDPOINT, res)
	return res
}

func OrbitShip(symbol string) ApiResult {
	endpoint := fmt.Sprintf(ORBIT_SHIP_ENDPOINT, symbol)
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetError(&errResp).
		Post(url + endpoint)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodPost, endpoint, res)
	return res
}

func DockShip(symbol string) ApiResult {
	endpoint := fmt.Sprintf(DOCK_SHIP_ENDPOINT, symbol)
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetError(&errResp).
		Post(url + endpoint)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodPost, endpoint, res)
	return res
}
