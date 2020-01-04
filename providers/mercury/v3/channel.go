package v3

// 2.1. ЗАПРОС НА ТЕСТИРОВАНИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelTest() error {
	request := &Request{
		Address: m.options.address,
		Code:    RequestCodeChannelTest,
	}

	resp, err := m.RequestRaw(request)

	if err != nil {
		return err
	}

	return ResponseError(request, resp)
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelOpen(level accessLevel, password LevelPassword) error {
	l := byte(level)

	request := &Request{
		Address:       m.options.address,
		Code:          RequestCodeChannelOpen,
		ParameterCode: &l,
		Parameters:    password.Bytes(),
	}

	resp, err := m.RequestRaw(request)

	if err != nil {
		return err
	}

	return ResponseError(request, resp)
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelClose() error {
	request := &Request{
		Address: m.options.address,
		Code:    RequestCodeChannelClose,
	}

	resp, err := m.RequestRaw(request)

	if err != nil {
		return err
	}

	return ResponseError(request, resp)
}
