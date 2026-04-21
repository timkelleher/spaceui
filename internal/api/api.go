package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/timkelleher/spaceui/internal/logger"
	resty "resty.dev/v3"
)

var (
	client *resty.Client

	url    = "https://api.spacetraders.io/v2/"
	apiKey = ""
)

func init() {
	client = resty.New()
}

func SetApiKey(val string) {
	apiKey = val
}

func Close() {
	if client != nil {
		client.Close()
	}
}

func logResponse(method, endpoint string, res ApiResult) {
	err := res.Err
	statusCode := 0
	if res.Resp != nil {
		statusCode = res.Resp.StatusCode()

		if res.Resp.IsError() {
			err = errors.New(res.ErrResp.Error.Message)
			statusCode = res.ErrResp.Error.Code
		}
	}

	switch statusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		logger.Info(fmt.Sprintf("%s %s [green]%d[-]", method, endpoint, statusCode))
	default:
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		logger.Warn(fmt.Sprintf("%s %s [red]%d[-] [red]%s[-]", method, endpoint, statusCode, errMsg))
	}
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Expected string `json:"expected"`
			Actual   string `json:"actual"`
		} `json:"data"`
		RequestID string `json:"requestId"`
	} `json:"error"`
}

type ListMetaResponse struct {
	Total int `json:"total,omitempty"`
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}
