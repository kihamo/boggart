package elektroset

import (
	"github.com/kihamo/boggart/providers/integratorit/internal"
)

const (
	TariffPlanEntityPay   = 1 // оплата
	TariffPlanEntityBill  = 2 // счет
	TariffPlanEntityValue = 3 // показания
)

type Provider struct {
	ID         uint64 `json:"kd_reg"`
	SchemaName string `json:"nm_schema"`
}

type House struct {
	ID       uint64 `json:"id_house"`
	Name     string `json:"nm_house"`
	Address  string `json:"nm_address"`
	Services []struct {
		ID              uint64 `json:"id_service"`
		ProviderID      uint64 `json:"kd_provider"`
		DeviceTypeID    uint64 `json:"kd_dervice_type"`
		StatusID        uint64 `json:"kd_status"`
		ProviderName    string `json:"nm_provider"`
		ServiceTypeName string `json:"nm_service_type"`
		TypeName        string `json:"nm_type"`
		AccountID       string `json:"nn_ls"`
		BalanceRAW      string `json:"vl_balance"`
		Balance         float64
		Provider        Provider `json:"vl_provider"`
	} `json:"services"`
}

type BalanceDetail struct {
	DatetimeEntity     internal.Time  `json:"dt_entity"`
	TariffPlanEntityID int64          `json:"kd_tp_entity"`
	DatetimePeriod     *internal.Time `json:"dt_period"`
	Sum                *float64       `json:"sm_entity"`
	StatusBillID       *int64         `json:"kd_st_bill"` // всегда 0 у выставленных счетов, у остальных null
	Zone1Name          *string        `json:"nm_zone1"`
	Zone2Name          *string        `json:"nm_zone2"`
	Zone3Name          *string        `json:"nm_zone3"`
	T1Value            *float64       `json:"vl_t1"`
	T2Value            *float64       `json:"vl_t2"`
	T3Value            *float64       `json:"vl_t3"`
	ConsumpValue       *string        `json:"vl_consump"`
}

type IndicationInfo struct {
	DatetimePreviousIndication *internal.Date `json:"dt_previous_ind"`
	ResultID                   uint64         `json:"kd_result"`
	ResultName                 interface{}    `json:"nm_result"`
	Zone1Name                  string         `json:"nm_zone1"`
	Zone2Name                  string         `json:"nm_zone2"`
	Zone3Name                  string         `json:"nm_zone3"`
	MeterName                  string         `json:"nn_meter"`
	Cons1CurrentDate           interface{}    `json:"vl_curdate_cons1"`
	Cons2CurrentDate           interface{}    `json:"vl_curdate_cons2"`
	Cons3CurrentDate           interface{}    `json:"vl_curdate_cons3"`
	T1CurrentDate              interface{}    `json:"vl_curdate_t1"`
	T2CurrentDate              interface{}    `json:"vl_curdate_t2"`
	T3CurrentDate              interface{}    `json:"vl_curdate_t3"`
	Cons1Previous              uint64         `json:"vl_previous_cons1"`
	Cons2Previous              uint64         `json:"vl_previous_cons2"`
	Cons3Previous              uint64         `json:"vl_previous_cons3"`
	T1Previous                 uint64         `json:"vl_previous_t1"`
	T2Previous                 uint64         `json:"vl_previous_t2"`
	T3Previous                 uint64         `json:"vl_previous_t3"`
	Quantity                   float64        `json:"vl_quantity"`
	MeterTypeName              string         `json:"vl_tp_meter"`
}

type Rate struct {
	DatetimeStart  *internal.Date `json:"dt_start_act"`
	TimeZone1Name  string         `json:"nm_time_zone1"`
	TimeZone2Name  string         `json:"nm_time_zone2"`
	TimeZone3Name  string         `json:"nm_time_zone3"`
	TariffPlanName string         `json:"nm_tp_rate"`
	Zone1Name      string         `json:"nm_zone1"`
	Zone2Name      string         `json:"nm_zone2"`
	Zone3Name      string         `json:"nm_zone3"`
	Rate1Value     float64        `json:"vl_rate_zone1"`
	Rate2Value     float64        `json:"vl_rate_zone2"`
	Rate3Value     float64        `json:"vl_rate_zone3"`
}

type BillFile struct {
	Link string `json:"link_bill"`
}

type Meter struct {
	CalibrationDate internal.Date `json:"dt_calibr"`
	PhaseName       string        `json:"nm_phase_meter"`
	State           string        `json:"nm_state_meter"`
	Model           string        `json:"nm_tp_meter"`
	Rate            string        `json:"nm_tp_rate"`
	Number          string        `json:"nn_meter"`
	Quantity        float64       `json:"vl_quantity"`
}
