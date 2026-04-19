package api

import (
	"fmt"
	"net/http"
	"strconv"
)

const GET_SYSTEMS_ENDPOINT = "/systems"

type SystemsResponse struct {
	Systems []System         `json:"data"`
	Meta    ListMetaResponse `json:"meta"`
}

type System struct {
	Symbol        string     `json:"symbol"`
	SectorSymbol  string     `json:"sectorSymbol"`
	Type          string     `json:"type"`
	X             int        `json:"x"`
	Y             int        `json:"y"`
	Waypoints     []Waypoint `json:"waypoints"`
	Factions      []Faction  `json:"factions"`
	Constellation string     `json:"constellation"`
	Name          string     `json:"name"`
}

type Waypoint struct {
	Symbol   string    `json:"symbol"`
	Type     string    `json:"type"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
	Orbitals []Orbital `json:"orbitals"`
}

type Orbital struct {
	Symbol string `json:"symbol"`
}

// TODO
type Faction struct {
	Symbol string `json:"symbol"`
}

func GetSystems(page int) (*SystemsResponse, ApiResult) {
	var obj SystemsResponse
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&obj).
		SetError(&errResp).
		SetQueryParams(map[string]string{"page": strconv.Itoa(page)}).
		Get(url + GET_SYSTEMS_ENDPOINT)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodGet, GET_SYSTEMS_ENDPOINT+fmt.Sprintf(" (%d)", page), res)
	return &obj, res
}
