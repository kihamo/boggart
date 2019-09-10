package elektroset

type Account struct {
	ID       uint64 `json:"id_house"`
	Number   string `json:"nm_house"`
	Services []struct {
		ID            uint64 `json:"id_service"`
		KDProvider    uint64 `json:"kd_provider"`
		KDDerviceType uint64 `json:"kd_dervice_type"`
		KDStatus      uint64 `json:"kd_status"`
		NMProvider    string `json:"nm_provider"`
		NMServiceType string `json:"nm_service_type"`
		NMType        string `json:"nm_type"`
		NNAccountRAW  string `json:"nn_ls"`
		NNAccount     uint64
		BalanceRAW    string `json:"vl_balance"`
		Balance       float64
		Provider      struct {
			KDReg    uint64 `json:"kd_reg"`
			NMSchema string `json:"nm_schema"`
		} `json:"vl_provider"`
	} `json:"services"`
}
