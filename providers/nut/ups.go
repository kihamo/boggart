package nut

type UPS struct {
	ListUPS

	session *Session
}

type Variable struct {
	ListVariable

	session *Session

	Type        Type
	Value       interface{}
	Description string
}

type Command struct {
	ListCommand

	session *Session

	Description string
}

func (c *Client) UPS() ([]UPS, error) {
	session := c.Session()

	listUPS, err := session.ListUPS()
	if err != nil {
		return nil, err
	}

	result := make([]UPS, len(listUPS))
	for i, item := range listUPS {
		result[i].ListUPS = item
		result[i].session = session
	}

	return result, err
}

func (u UPS) Variables() ([]Variable, error) {
	listVariables, err := u.session.ListVariables(u.Name)
	if err != nil {
		return nil, err
	}

	result := make([]Variable, len(listVariables))
	for i, item := range listVariables {
		result[i].ListVariable = item
		result[i].session = u.session

		result[i].Type, err = u.session.Type(u.Name, item.Name)
		if err != nil {
			return nil, err
		}
		result[i].Value = result[i].Type.ConvertValue(item.Value)

		result[i].Description, err = u.session.Description(u.Name, item.Name)
		if err != nil {
			return nil, err
		}
	}

	return result, err
}

func (u UPS) Commands() ([]Command, error) {
	listCommands, err := u.session.ListCommands(u.Name)
	if err != nil {
		return nil, err
	}

	result := make([]Command, len(listCommands))
	for i, item := range listCommands {
		result[i].ListCommand = item
		result[i].session = u.session

		result[i].Description, err = u.session.CommandDescription(u.Name, item.Name)
		if err != nil {
			return nil, err
		}
	}

	return result, err
}

func (v Variable) Set(value interface{}) error {
	return v.session.Set(v.UPS, v.Name, value)
}

func (c Command) Call() error {
	return c.session.InstantCommand(c.UPS, c.Name)
}
