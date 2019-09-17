package elektroset

import (
	"github.com/kihamo/boggart/providers/integratorit/internal"
)

type Provider struct {
	KDReg    uint64 `json:"kd_reg"`
	NMSchema string `json:"nm_schema"`
}

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
		Provider      Provider `json:"vl_provider"`
	} `json:"services"`
}

type BalanceDetail struct {
	DatetimePeriod *internal.Time `json:"dt_period"`
	DatetimeEntity *internal.Time `json:"dt_entity"`
	Zone1          *string        `json:"nm_zone1"`
	Zone2          *string        `json:"nm_zone2"`
	Zone3          *string        `json:"nm_zone3"`
	ValueT1        *float64       `json:"vl_t1"`
	ValueT2        *float64       `json:"vl_t2"`
	ValueT3        *float64       `json:"vl_t3"`
}
