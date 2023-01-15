package internal

type Response struct {
	ErrorCode int         `json:"error_code"`
	Result    interface{} `json:"result"`
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) WithResult(result interface{}) *Response {
	r.Result = result
	return r
}
