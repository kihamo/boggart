package mosenergosbyt

import "time"

const (
	OwnerTypeOwner      uint8 = iota + 1 // Собственник
	OwnerTypeOther                       // Другое
	OwnerTypeTenant                      // Наниматель
	OwnerTypeRegistered                  // Зарегистрированный
	OwnerTypeResides                     // Проживает
)

// https://my.mosenergosbyt.ru/static/js/constants/ProvidersConfig.js
const (
	ProviderIDMosEnergoSbyt   uint8 = iota + 1 // MES_KD_PROVIDER - Мосэнергосбыт
	ProviderIDMosOblERC                        // MOE_KD_PROVIDER
	ProviderIDTomskEnergoSbyt                  // TMK_NRG_KD_PROVIDER - Томскэнергосбыт
	ProviderIDTomskRTS                         // TMK_RTS_KD_PROVIDER
	ProviderIDUfa                              // UFA_KD_PROVIDER - Уфа
	ProviderIDTKO                              // TKO_KD_PROVIDER
	ProviderIDVLG                              // VLG_KD_PROVIDER
	ProviderIDOrelNGR                          // ORL_KD_PROVIDER
	ProviderIDOrelEPD                          // ORL_EPD_KD_PROVIDER - Орел ЕПД
	ProviderIDAlt                              // ALT_KD_PROVIDER
	ProviderIDTmb                              // TMB_KD_PROVIDER
	ProviderIDVld                              // VLD_KD_PROVIDER
	ProviderIDSar                              // SAR_KD_PROVIDER
	ProviderIDKsg                              // KSG_KD_PROVIDER
)

// https://my.mosenergosbyt.ru/static/js/constants/types/AccountTypes.js

type Account struct {
	Data struct {
		OwnerType   uint8  `json:"kd_ls_owner_type"`
		IDTU        uint64 `json:"id_tu"`
		Street      string `json:"nm_street"`
		DisplayName string `json:"nn_ls_disp"`
		SNILS       string `json:"NN_SNILS"`
		Reg         uint64 `json:"kd_reg"`
	}
	ServiceID     uint64 `json:"id_service"`
	ServiceType   uint64 `json:"kd_service_type"`
	StatusID      uint64 `json:"kd_status"`
	LockMessage   string `json:"nm_lock_msg"`
	Description   string `json:"nm_ls_description"`
	TypeName      string `json:"nm_type"`
	AccountID     string `json:"nn_ls"`
	Group         string `json:"nm_ls_group"`
	GroupFull     string `json:"nm_ls_group_full"`
	GroupEditable bool   `json:"pr_ls_group_edit"`
	ProviderID    uint8  `json:"kd_provider"`
	ProviderName  string `json:"nm_provider"`
	ProviderRAW   string `json:"vl_provider"`
	Provider      map[string]interface{}
}

type Balance struct {
	Total    float64
	Services []ServiceBalance
}

type ServiceBalance struct {
	ID    string
	Name  string
	Total float64
}

type Bill struct {
	ID     string
	Period time.Time
	Amount float64
}
