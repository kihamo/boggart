package miio

type Response struct {
	ID     uint32      `json:"id"`
	Result interface{} `json:"result"`
}

type Request struct {
	ID     uint32      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}
