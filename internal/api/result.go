package api

import resty "resty.dev/v3"

type ApiResult struct {
	Err     error
	ErrResp ErrorResponse
	Resp    *resty.Response
}

func (ar ApiResult) Error() bool {
	return ar.Resp.IsError() || ar.Err != nil
}

func (ar ApiResult) ErrorMessage() string {
	if !ar.Error() {
		return ""
	}
	if ar.ErrResp.Error.Message != "" {
		return ar.ErrResp.Error.Message
	}
	if ar.Err != nil {
		return ar.Err.Error()
	}
	return "unknown_error"
}
