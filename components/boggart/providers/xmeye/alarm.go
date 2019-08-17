package xmeye

func (c *Client) AlarmStart() error {
	_, err := c.Call(CmdGuardRequest, nil)
	return err
}

func (c *Client) AlarmStop() error {
	_, err := c.Call(CmdUnGuardRequest, nil)
	return err
}

func (c *Client) AlarmInfo() (*AlarmInfo, error) {
	response := &AlarmInfo{}

	err := c.CallWithResult(CmdAlarmRequest, nil, response)

	return response, err
}
