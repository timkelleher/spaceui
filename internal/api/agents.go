package api

import "net/http"

const GET_AGENT_ENDPOINT = "/my/agent"

type AgentResponse struct {
	Agent Agent `json:"data,omitempty"`
}

type Agent struct {
	AccountID       string `json:"accountId,omitempty"`
	Symbol          string `json:"symbol,omitempty"`
	Headquarters    string `json:"headquarters,omitempty"`
	Credits         int    `json:"credits,omitempty"`
	StartingFaction string `json:"startingFaction,omitempty"`
	ShipCount       int    `json:"shipCount,omitempty"`
}

func GetAgent() (*AgentResponse, ApiResult) {
	var agent AgentResponse
	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&agent).
		Get(url + GET_AGENT_ENDPOINT)

	logResponse(http.MethodGet, GET_AGENT_ENDPOINT, res.StatusCode(), err)
	return &agent, ApiResult{Resp: res, Err: err}
}
