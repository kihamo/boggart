package internal

type response struct {
	Data     interface{} `json:"data"`
	MetaData struct {
		ResponseTime float64 `json:"responseTime"`
	} `json:"metaData"`
	Success      bool   `json:"success"`
	Total        uint64 `json:"total"`
	ErrorCode    uint64 `json:"err_code"`
	ErrorMessage string `json:"err_text"`
	ErrorID      string `json:"err_id"`
}
