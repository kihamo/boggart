package v3

// 2.1. ЗАПРОС НА ТЕСТИРОВАНИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelTest() error {
	response, err := m.InvokeRaw(NewRequest(CodeChannelTest))
	if err != nil {
		return err
	}

	return response.GetError()
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelOpen(level uint8, password LevelPassword) error {
	request := NewRequest(CodeChannelOpen).
		WithParameterCode(level).
		WithParameters(password.Bytes())

	response, err := m.InvokeRaw(request)
	if err != nil {
		return err
	}

	return response.GetError()
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (m *MercuryV3) ChannelClose() error {
	response, err := m.InvokeRaw(NewRequest(CodeChannelClose))
	if err != nil {
		return err
	}

	return response.GetError()
}
