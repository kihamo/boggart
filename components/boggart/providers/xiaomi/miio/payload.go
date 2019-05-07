package miio

import (
	"strings"
)

type Response struct {
	ID     uint32      `json:"id"`
	Result interface{} `json:"result"`
}

type ResponseOK struct {
	Response

	Result []string `json:"result"`
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
