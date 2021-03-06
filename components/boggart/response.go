package boggart

type ResponseJSON struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func NewResponseJSON() *ResponseJSON {
	return &ResponseJSON{}
}

func (r *ResponseJSON) Failed(message string) *ResponseJSON {
	r.Result = "failed"
	r.Message = message

	return r
}

func (r *ResponseJSON) FailedError(err error) *ResponseJSON {
	return r.Failed(err.Error())
}

func (r *ResponseJSON) Success(message string) *ResponseJSON {
	r.Result = "success"
	r.Message = message

	return r
}

type ResponseDataTable struct {
	Draw     int         `json:"draw"`
	Total    int         `json:"recordsTotal"`
	Filtered int         `json:"recordsFiltered"`
	Data     interface{} `json:"data"`
	Error    string      `json:"error,omitempty"`
}

func NewResponseDataTable() *ResponseDataTable {
	return &ResponseDataTable{
		Draw: 1,
	}
}
