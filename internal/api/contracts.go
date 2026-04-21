package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/timkelleher/spaceui/internal/logger"
)

const (
	GET_CONTRACTS_ENDPOINT      = "/my/contracts"
	NEGOTIATE_CONTRACT_ENDPOINT = "/my/ships/%s/negotiate/contract"
	ACCEPT_CONTRACT_ENDPOINT    = "/my/contracts/%s/accept"
)

type ContractsResponse struct {
	Contracts []Contract       `json:"data,omitempty"`
	Meta      ListMetaResponse `json:"meta"`
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
	var obj ContractsResponse
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&obj).
		SetError(&errResp).
		Get(url + GET_CONTRACTS_ENDPOINT)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodGet, GET_CONTRACTS_ENDPOINT, res)
	return &obj, res
}

type NegotiateContractResponse struct {
	Data struct {
		Contract Contract `json:"contract,omitempty"`
	} `json:"data"`
}

func NegotiateContract(shipSymbol string) (*NegotiateContractResponse, ApiResult) {
	endpoint := fmt.Sprintf(NEGOTIATE_CONTRACT_ENDPOINT, shipSymbol)

	var obj NegotiateContractResponse
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&obj).
		SetError(&errResp).
		Post(url + endpoint)

	logger.Info(fmt.Sprintf("%d", resp.StatusCode()))
	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodPost, endpoint, res)
	return &obj, res
}

type AcceptContractResponse struct {
	Data struct {
		Contract Contract `json:"contract,omitempty"`
		Agent    Agent    `json:"agent,omitempty"`
	} `json:"data"`
}

func AcceptContract(id string) (*AcceptContractResponse, ApiResult) {
	endpoint := fmt.Sprintf(ACCEPT_CONTRACT_ENDPOINT, id)

	var obj AcceptContractResponse
	var errResp ErrorResponse

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(apiKey).
		SetResult(&obj).
		SetError(&errResp).
		Post(url + endpoint)

	res := ApiResult{Resp: resp, ErrResp: errResp, Err: err}
	logResponse(http.MethodPost, endpoint, res)
	return &obj, res
}
