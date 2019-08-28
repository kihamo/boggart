package v3

// 2.1. ЗАПРОС НА ТЕСТИРОВАНИЕ КАНАЛА СВЯЗИ
func (d *MercuryV3) ChannelTest() error {
	resp, err := d.Request(&Request{
		Address: d.address,
		Code:    RequestCodeChannelTest,
	})

	if err != nil {
		return err
	}

	return ResponseError(resp)
}

// 2.2. ЗАПРОСЫ НА ОТКРЫТИЕ/ЗАКРЫТИЕ КАНАЛА СВЯЗИ
func (d *MercuryV3) ChannelOpen(level accessLevel, password LevelPassword) error {
	l := byte(level)

	resp, err := d.Request(&Request{
		Address:       d.address,
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
func (d *MercuryV3) ChannelClose() error {
	resp, err := d.Request(&Request{
		Address: d.address,
		Code:    RequestCodeChannelClose,
	})

	if err != nil {
		return err
	}

	return ResponseError(resp)
}
