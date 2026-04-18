package api

import (
	"net/http"
	"time"
)

const GET_CONTRACTS_ENDPOINT = "/my/contracts"

type ContractsResponse struct {
	Contracts []Contract   `json:"data,omitempty"`
	Meta      ContractMeta `json:"meta,omitempty"`
}

type Contract struct {
	ID               string    `json:"id,omitempty"`
	FactionSymbol    string    `json:"factionSymbol,omitempty"`
	Type             string    `json:"type,omitempty"`
	Terms            Terms     `json:"terms,omitempty"`
	Accepted         bool      `json:"accepted,omitempty"`
	Fulfilled        bool      `json:"fulfilled,omitempty"`
	Expiration       time.Time `json:"expiration,omitempty"`
	DeadlineToAccept time.Time `json:"deadlineToAccept,omitempty"`
}

func (c Contract) HasAccepted() string {
	if c.Accepted {
		return "[green]Accepted[-]"
	}
	return "Not Accepted"
}

func (c Contract) HasFulfilled() string {
	if c.Fulfilled {
		return "[green]Fulfilled[-]"
	}
	return "Not Fulfilled"
}

type ContractMeta struct {
	Total int `json:"total,omitempty"`
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

type Terms struct {
	Deadline time.Time `json:"deadline,omitempty"`
	Payment  Payment   `json:"payment,omitempty"`
	Deliver  []Deliver `json:"deliver,omitempty"`
}

type Payment struct {
	OnAccepted  int `json:"onAccepted,omitempty"`
	OnFulfilled int `json:"onFulfilled,omitempty"`
}

type Deliver struct {
	TradeSymbol       string `json:"tradeSymbol,omitempty"`
	DestinationSymbol string `json:"destinationSymbol,omitempty"`
	UnitsRequired     int    `json:"unitsRequired,omitempty"`
	UnitsFulfilled    int    `json:"unitsFulfilled,omitempty"`
}

func GetContracts() (*ContractsResponse, ApiResult) {
	if client == nil {
		Init()
	}

	var contracts ContractsResponse
	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&contracts).
		Get(url + GET_CONTRACTS_ENDPOINT)

	logResponse(http.MethodGet, GET_CONTRACTS_ENDPOINT, res.StatusCode(), err)
	return &contracts, ApiResult{Resp: res, Err: err}
}
