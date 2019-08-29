package v3

// 2.1. ЗАПРОС НА ТЕСТИРОВАНИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelTest() error {
	resp, err := m.Request(&Request{
		Address: m.address,
		Code:    RequestCodeChannelTest,
	})

	if err != nil {
		return err
	}

	return ResponseError(resp)
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelOpen(level accessLevel, password LevelPassword) error {
	l := byte(level)

	resp, err := m.Request(&Request{
		Address:       m.address,
		Code:          RequestCodeChannelOpen,
		ParameterCode: &l,
		Parameters:    password.Bytes(),
	})

	if err != nil {
		return err
	}

	return ResponseError(resp)
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelClose() error {
	resp, err := m.Request(&Request{
		Address: m.address,
		Code:    RequestCodeChannelClose,
	})

	if err != nil {
		return err
	}

	return ResponseError(resp)
}
