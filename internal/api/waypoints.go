package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const GET_WAYPOINTS_ENDPOINT = "/systems/%s/waypoints"

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
