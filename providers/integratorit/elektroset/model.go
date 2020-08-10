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
		KDDeviceType  uint64 `json:"kd_dervice_type"`
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

type IndicationInfo struct {
	DatetimePreviousIndication *internal.Date `json:"dt_previous_ind"`
	KDResult                   uint64         `json:"kd_result"`
	NMResult                   interface{}    `json:"nm_result"`
	Zone1                      string         `json:"nm_zone1"`
	Zone2                      string         `json:"nm_zone2"`
	Zone3                      string         `json:"nm_zone3"`
	NNMeter                    string         `json:"nn_meter"`
	CurrentDateCons1           interface{}    `json:"vl_curdate_cons1"`
	CurrentDateCons2           interface{}    `json:"vl_curdate_cons2"`
	CurrentDateCons3           interface{}    `json:"vl_curdate_cons3"`
	CurrentDateT1              interface{}    `json:"vl_curdate_t1"`
	CurrentDateT2              interface{}    `json:"vl_curdate_t2"`
	CurrentDateT3              interface{}    `json:"vl_curdate_t3"`
	PreviousCons1              uint64         `json:"vl_previous_cons1"`
	PreviousCons2              uint64         `json:"vl_previous_cons2"`
	PreviousCons3              uint64         `json:"vl_previous_cons3"`
	PreviousT1                 uint64         `json:"vl_previous_t1"`
	PreviousT2                 uint64         `json:"vl_previous_t2"`
	PreviousT3                 uint64         `json:"vl_previous_t3"`
	Quantity                   float64        `json:"vl_quantity"`
	TariffPlanMeter            string         `json:"vl_tp_meter"`
}

type Rate struct {
	DatetimeStart *internal.Date `json:"dt_start_act"`
	Time1         string         `json:"nm_time_zone1"`
	Time2         string         `json:"nm_time_zone2"`
	Time3         string         `json:"nm_time_zone3"`
	TariffPlan    string         `json:"nm_tp_rate"`
	Zone1         string         `json:"nm_zone1"`
	Zone2         string         `json:"nm_zone2"`
	Zone3         string         `json:"nm_zone3"`
	Rate1         float64        `json:"vl_rate_zone1"`
	Rate2         float64        `json:"vl_rate_zone2"`
	Rate3         float64        `json:"vl_rate_zone3"`
}
