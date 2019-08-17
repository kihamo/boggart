package xmeye

import (
	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

func (c *Client) SystemInfo() (*internal.SystemInfo, error) {
	var result struct {
		internal.Response
		SystemInfo internal.SystemInfo
	}

	err := c.CmdWithResult(CmdSystemInfoRequest, "SystemInfo", &result)
	if err != nil {
		return nil, err
	}

	return &result.SystemInfo, err
}

func (c *Client) OEMInfo() (*internal.OEMInfo, error) {
	var result struct {
		internal.Response
		OEMInfo internal.OEMInfo
	}

	err := c.CmdWithResult(CmdSystemInfoRequest, "OEMInfo", &result)
	if err != nil {
		return nil, err
	}

	return &result.OEMInfo, err
}

func (c *Client) StorageInfo() (*internal.StorageInfo, error) {
	var result struct {
		internal.Response
		StorageInfo internal.StorageInfo
	}

	err := c.CmdWithResult(CmdSystemInfoRequest, "StorageInfo", &result)
	if err != nil {
		return nil, err
	}

	return &result.StorageInfo, err
}

func (c *Client) WorkState() (*internal.WorkState, error) {
	var result struct {
		internal.Response
		WorkState internal.WorkState
	}

	err := c.CmdWithResult(CmdSystemInfoRequest, "WorkState", &result)
	if err != nil {
		return nil, err
	}

	return &result.WorkState, err
}
