package api

import (
	"fmt"
	"net/http"
	"time"
)

const (
	GET_CONTRACTS_ENDPOINT      = "/my/contracts"
	NEGOTIATE_CONTRACT_ENDPOINT = "/my/ships/%s/negotiate/contract"
	ACCEPT_CONTRACT_ENDPOINT    = "/my/contracts/%s/accept"
)

type ContractsResponse struct {
	Contracts []Contract   `json:"data,omitempty"`
	Meta      ContractMeta `json:"meta,omitempty"`
}

type Contract struct {
	ID            string `json:"id"`
	FactionSymbol string `json:"factionSymbol"`
	Type          string `json:"type"`
	Terms         struct {
		Deadline time.Time `json:"deadline"`
		Payment  struct {
			OnAccepted  int `json:"onAccepted"`
			OnFulfilled int `json:"onFulfilled"`
		} `json:"payment"`
		Deliver []struct {
			TradeSymbol       string `json:"tradeSymbol"`
			DestinationSymbol string `json:"destinationSymbol"`
			UnitsRequired     int    `json:"unitsRequired"`
			UnitsFulfilled    int    `json:"unitsFulfilled"`
		} `json:"deliver"`
	} `json:"terms"`
	Accepted         bool      `json:"accepted"`
	Fulfilled        bool      `json:"fulfilled"`
	Expiration       time.Time `json:"expiration"`
	DeadlineToAccept time.Time `json:"deadlineToAccept"`
}

func (c Contract) State() string {
	if !c.Accepted {
		return "Not Accepted"
	} else if !c.Fulfilled {
		return "Accepted"
	} else {
		return "Fulfilled"
	}
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
	var res ContractsResponse
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&res).
		Get(url + GET_CONTRACTS_ENDPOINT)

	logResponse(http.MethodGet, GET_CONTRACTS_ENDPOINT, resp.StatusCode(), err)
	return &res, ApiResult{Resp: resp, Err: err}
}

type NegotiateContractResponse struct {
	Data struct {
		Contract Contract `json:"contract,omitempty"`
	} `json:"data"`
}

func NegotiateContract(shipSymbol string) (*NegotiateContractResponse, ApiResult) {
	endpoint := fmt.Sprintf(NEGOTIATE_CONTRACT_ENDPOINT, shipSymbol)

	var res NegotiateContractResponse
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&res).
		Post(url + endpoint)

	logResponse(http.MethodPost, endpoint, resp.StatusCode(), err)
	return &res, ApiResult{Resp: resp, Err: err}
}

type AcceptContractResponse struct {
	Data struct {
		Contract Contract `json:"contract,omitempty"`
		Agent    Agent    `json:"agent,omitempty"`
	} `json:"data"`
}

func AcceptContract(id string) (*AcceptContractResponse, ApiResult) {
	endpoint := fmt.Sprintf(ACCEPT_CONTRACT_ENDPOINT, id)

	var res AcceptContractResponse
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&res).
		Post(url + endpoint)

	logResponse(http.MethodPost, endpoint, resp.StatusCode(), err)
	return &res, ApiResult{Resp: resp, Err: err}
}
