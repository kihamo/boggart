package pulsar_heat_meter

type Config struct {
	RS485 struct {
		Address string `valid:"required"`
		Timeout string
	} `valid:"required"`
	Address      string
	Input1Offset float64 `mapstructure:"input1_offset",valid:"float"`
	Input2Offset float64 `mapstructure:"input2_offset",valid:"float"`
	Input3Offset float64 `mapstructure:"input3_offset",valid:"float"`
	Input4Offset float64 `mapstructure:"input4_offset",valid:"float"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
