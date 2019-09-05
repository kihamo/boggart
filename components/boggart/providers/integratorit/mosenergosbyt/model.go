package mosenergosbyt

import (
	"github.com/kihamo/boggart/components/boggart/providers/integratorit/internal"
)

type Account struct {
	Data struct {
		KDAccountOwnerType uint64 `json:"kd_ls_owner_type"`
		IDTU               uint64 `json:"id_tu"`
		NMStreet           string `json:"nm_street"`
		NNAccountDisp      string `json:"nn_ls_disp"`
	}
	IDService          uint64 `json:"id_service"`
	KDProvider         uint64 `json:"kd_provider"`
	KDServiceType      uint64 `json:"kd_service_type"`
	KDStatus           uint64 `json:"kd_status"`
	NKLockMessage      string `json:"nm_lock_msg"`
	NMAccountGroup     string `json:"nm_ls_group"`
	NMAccountGroupFull string `json:"nm_ls_group_full"`
	NMProvider         string `json:"nm_provider"`
	NMType             string `json:"nm_type"`
	NNAccount          string `json:"nn_ls"`
	PRAccountGroupEdit bool   `json:"pr_ls_group_edit"`
	ProviderRAW        string `json:"vl_provider"`
	Provider           struct {
		IDAbonent uint64 `json:"id_abonent"`
	}
}

type Balance struct {
	DeptPeriodBalance string  `json:"dt_period_balance"`
	Advance           float64 `json:"sm_advance"`
	Balance           float64 `json:"sm_balance"`
	Charged           float64 `json:"sm_charged"`
	Insurance         float64 `json:"sm_insurance"`
	Payed             float64 `json:"sm_payed"`
	Penalty           float64 `json:"sm_penalty"`
	StartAdvance      float64 `json:"sm_start_advance"`
	StartDebt         float64 `json:"sm_start_debt"`
	Tovkgo            float64 `json:"sm_tovkgo"`
	Services          []struct {
		Service      string  `json:"nm_service"`
		Advance      float64 `json:"sm_advance"`
		Balance      float64 `json:"sm_balance"`
		Charged      float64 `json:"sm_charged"`
		Payed        float64 `json:"sm_payed"`
		Penalty      float64 `json:"sm_penalty"`
		StartAdvance float64 `json:"sm_start_advance"`
		StartDebt    float64 `json:"sm_start_debt"`
	} `json:"child"`
}

type Payment struct {
	Date   internal.Time `json:"dt_pay"`
	Agent  string        `json:"nm_agnt"`
	State  string        `json:"nm_pay_state"`
	Amount float64       `json:"sm_pay"`
}

type Charge struct {
	Create   internal.Time `json:"dt_create"`
	Period   internal.Date `json:"dt_period"`
	Services []struct {
		KDType         uint64  `json:"kd_child_type"`
		Benefits       float64 `json:"sm_benefits"`
		Charged        float64 `json:"sm_charged"`
		Insurance      uint64  `json:"sm_insurance"`
		Payed          float64 `json:"sm_payed"`
		Penalty        float64 `json:"sm_penalty"`
		Recalculations float64 `json:"sm_recalculations"`
		Start          float64 `json:"sm_start"`
		Total          float64 `json:"sm_total"`
		Tovkgo         float64 `json:"sm_tovkgo"`
		ReportUUID     string  `json:"vl_report_uuid"`
		List           []struct {
			MeasureUnit    string  `json:"nm_measure_unit"`
			Service        string  `json:"nm_service"`
			Benefits       float64 `json:"sm_benefits"`
			Charged        float64 `json:"sm_charged"`
			Payed          float64 `json:"sm_payed"`
			Penalty        float64 `json:"sm_penalty"`
			Recalculations float64 `json:"sm_recalculations"`
			Start          float64 `json:"sm_start"`
			Total          float64 `json:"sm_total"`
			ChargedVolume  float64 `json:"vl_charged_volume"`
			Tariff         float64 `json:"vl_tariff"`
		} `json:"child"`
	} `json:"child"`
}
