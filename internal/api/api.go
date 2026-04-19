package api

import (
	"fmt"

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

type ApiResult struct {
	Err  error
	Resp *resty.Response
}

func SetApiKey(val string) {
	apiKey = val
}

func Close() {
	if client != nil {
		client.Close()
	}
}

func logResponse(method, endpoint string, statusCode int, err error) {
	if err != nil {
		logger.Error(fmt.Sprintf("%s %s [yellow]%d[-] [red]%s[-]", method, endpoint, statusCode, err.Error()))
	} else {
		logger.Info(fmt.Sprintf("%s %s [yellow]%d[-]", method, endpoint, statusCode))
	}
}
