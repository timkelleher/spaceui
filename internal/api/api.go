package api

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/timkelleher/spaceui/internal/logger"
	resty "resty.dev/v3"
)

var client *resty.Client
var url = "https://api.spacetraders.io/v2/"
var apiKey = ""

type ApiResult struct {
	Err  error
	Resp *resty.Response
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	apiKey = os.Getenv("SPACE_TRADERS_API_KEY")
	client = resty.New()
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
