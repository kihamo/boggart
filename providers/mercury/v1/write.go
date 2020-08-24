package v1

// true  / разрешает индикацию 1 тарифа
// true  / разрешает индикацию 2 тарифа
// false / разрешает индикацию 3 тарифа
// false / разрешает индикацию 4 тарифа
// true  / разрешает индикацию суммы
// false / разрешает индикацию мощности
// false / разрешает индикацию времени
// false / разрешает индикацию даты
func (m *MercuryV1) SetDisplayMode(mode *DisplayMode) error {
	request := NewPacket().
		WithCommand(CommandWriteDisplayMode).
		WithPayload([]byte{mode.Bit()})

	_, err := m.Invoke(request)

	return err
}

// default 10 10 5 30
// t1 / 10 / время индикации энергии не текущих тарифов и суммы
// t2 / 10 / время индикации энергии текущего тарифа
// t3 /  5 / время индикации мощности, времени и даты
// t4 / 30 / время индикации после нажатия кнопки
func (m *MercuryV1) SetDisplayTime(values *TariffValues) error {
	request := NewPacket().
		WithCommand(CommandWriteDisplayTime).
		WithPayload([]byte{uint8(values.Tariff1()), uint8(values.Tariff2()), uint8(values.Tariff3()), uint8(values.Tariff4())})

	_, err := m.Invoke(request)

	return err
}
