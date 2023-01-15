package internal

type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) WithMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) WithParams(params interface{}) *Request {
	r.Params = params
	return r
}
