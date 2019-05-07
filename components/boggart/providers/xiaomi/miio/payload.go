package miio

import (
	"strings"
)

type Response struct {
	ID     uint32      `json:"id"`
	Result interface{} `json:"result"`
}

type ResponseOK struct {
	ID     uint32   `json:"id"`
	Result []string `json:"result"`
}

type ResponseUnknownMethod struct {
	ID     uint32 `json:"id"`
	Result string `json:"result"`
}

type ResponseError struct {
	ID    uint32 `json:"id"`
	Error struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Request struct {
	ID     uint32      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func ResponseIsOK(response ResponseOK) bool {
	if len(response.Result) > 0 {
		return strings.EqualFold(response.Result[0], "ok")
	}

	return false
}
