package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	GET_WAYPOINTS_ENDPOINT         = "/systems/%s/waypoints"
	GET_WAYPOINT_MARKET_ENDPOINTS  = "/systems/%s/waypoints/%s/market"
	GET_WAYPOINT_SHIPYARD_ENDPOINT = "/systems/%s/waypoints/%s/shipyard"
)

// /////////////////////////////////////
// Waypoints
// /////////////////////////////////////

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

func (w Waypoint) IsMarketplace() bool {
	return w.hasTrait("Marketplace")
}

func (w Waypoint) IsShipyard() bool {
	return w.hasTrait("Shipyard")
}

func (w Waypoint) hasTrait(desired string) bool {
	for _, trait := range w.Traits {
		if trait.Name == desired {
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

// /////////////////////////////////////
// Waypoint Market
// /////////////////////////////////////

type WaypoiontMarketResponse struct {
	Marketplace `json:"data"`
}

type Marketplace struct {
	Symbol       string                   `json:"symbol"`
	Exports      []interface{}            `json:"exports"`
	Imports      []interface{}            `json:"imports"`
	Exchange     []MarketplaceExchange    `json:"exchange"`
	Transactions []MarketplaceTransaction `json:"transactions"`
	TradeGoods   []MarketplaceTradeGood   `json:"tradeGoods"`
}

type MarketplaceExchange struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MarketplaceTransaction struct {
	WaypointSymbol string    `json:"waypointSymbol"`
	ShipSymbol     string    `json:"shipSymbol"`
	TradeSymbol    string    `json:"tradeSymbol"`
	Type           string    `json:"type"`
	Units          int       `json:"units"`
	PricePerUnit   int       `json:"pricePerUnit"`
	TotalPrice     int       `json:"totalPrice"`
	Timestamp      time.Time `json:"timestamp"`
}

type MarketplaceTradeGood struct {
	Symbol        string `json:"symbol"`
	Type          string `json:"type"`
	TradeVolume   int    `json:"tradeVolume"`
	Supply        string `json:"supply"`
	PurchasePrice int    `json:"purchasePrice"`
	SellPrice     int    `json:"sellPrice"`
}

func GetMarketplace(waypoint Waypoint) (*WaypoiontMarketResponse, ApiResult) {
	endpoint := fmt.Sprintf(GET_WAYPOINT_MARKET_ENDPOINTS, waypoint.SystemSymbol, waypoint.Symbol)

	var obj WaypoiontMarketResponse
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

// /////////////////////////////////////
// Waypoint Shipyard
// /////////////////////////////////////

type WaypointShipyardResponse struct {
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

func GetShipyard(waypoint Waypoint) (*WaypointShipyardResponse, ApiResult) {
	endpoint := fmt.Sprintf(GET_WAYPOINT_SHIPYARD_ENDPOINT, waypoint.SystemSymbol, waypoint.Symbol)

	var obj WaypointShipyardResponse
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
