package xmeye

import (
	"time"
)

func (c *Client) OPTime() (*time.Time, error) {
	response := &OPTimeQuery{}

	err := c.CmdWithResult(CmdTimeRequest, "OPTimeQuery", response)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse("2006-01-02 15:04:05", response.OPTimeQuery)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *Client) SystemInfo() (*SystemInfo, error) {
	result := &SystemInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "SystemInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}
