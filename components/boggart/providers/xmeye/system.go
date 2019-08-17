package xmeye

import (
	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

func (c *Client) SystemFunctions() (*internal.SystemFunctions, error) {
	result := &internal.SystemFunctions{}

	err := c.CmdWithResult(CmdAbilityGetRequest, "SystemFunction", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) SystemInfo() (*internal.SystemInfo, error) {
	result := &internal.SystemInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "SystemInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) OEMInfo() (*internal.OEMInfo, error) {
	result := &internal.OEMInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "OEMInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) StorageInfo() (*internal.StorageInfo, error) {
	result := &internal.StorageInfo{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "StorageInfo", result)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (c *Client) WorkState() (*internal.WorkState, error) {
	result := &internal.WorkState{}

	err := c.CmdWithResult(CmdSystemInfoRequest, "WorkState", result)
	if err != nil {
		return nil, err
	}

	return result, err
}
