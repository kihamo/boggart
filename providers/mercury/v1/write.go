package v1

// true  / разрешает индикацию 1 тарифа
// true  / разрешает индикацию 2 тарифа
// false / разрешает индикацию 3 тарифа
// false / разрешает индикацию 4 тарифа
// true  / разрешает индикацию суммы
// false / разрешает индикацию мощности
// false / разрешает индикацию времени
// false / разрешает индикацию даты
func (m *MercuryV1) SetDisplayMode(t1, t2, t3, t4, amount, power, time, date bool) error {
	bit := 0

	if t1 {
		bit |= displayModeTariff1
	}

	if t2 {
		bit |= displayModeTariff2
	}

	if t3 {
		bit |= displayModeTariff3
	}

	if t4 {
		bit |= displayModeTariff4
	}

	if amount {
		bit |= displayModeAmount
	}

	if power {
		bit |= displayModePower
	}

	if time {
		bit |= displayModeTime
	}

	if date {
		bit |= displayModeDate
	}

	_, err := m.Request(&Request{
		Command: RequestCommandWriteDisplayMode,
		Payload: []byte{byte(bit)},
	})

	return err
}

// default 10 10 5 30
// t1 / 10 / время индикации энергии не текущих тарифов и суммы
// t2 / 10 / время индикации энергии текущего тарифа
// t3 /  5 / время индикации мощности, времени и даты
// t4 / 30 / время индикации после нажатия кнопки
func (m *MercuryV1) SetDisplayTime(t1, t2, t3, t4 uint64) error {
	_, err := m.Request(&Request{
		Command: RequestCommandWriteDisplayTime,
		Payload: []byte{byte(t1), byte(t2), byte(t3), byte(t4)},
	})

	return err
}
