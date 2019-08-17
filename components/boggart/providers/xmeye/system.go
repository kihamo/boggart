package xmeye

func (c *Client) SystemFunctions() (*SystemFunctions, error) {
	result := &SystemFunctions{}

	err := c.CmdWithResult(CmdAbilityGetRequest, "SystemFunction", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) SystemInfo() (*SystemInfo, error) {
	result := &SystemInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "SystemInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) OEMInfo() (*OEMInfo, error) {
	result := &OEMInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "OEMInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) StorageInfo() (*StorageInfo, error) {
	result := &StorageInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "StorageInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) WorkState() (*WorkState, error) {
	result := &WorkState{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "WorkState", result)
	if err != nil {
		return nil, err
	}

	return result, err
}
