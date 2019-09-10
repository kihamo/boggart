package xmeye

import (
	"context"
)

func (c *Client) SystemInfo(ctx context.Context) (*SystemInfo, error) {
	var result struct {
		Response
		SystemInfo SystemInfo
	}

	err := c.CmdWithResult(ctx, CmdSystemInfoRequest, "SystemInfo", &result)
	if err != nil {
		return nil, err
	}

	return &result.SystemInfo, err
}

func (c *Client) OEMInfo(ctx context.Context) (*OEMInfo, error) {
	var result struct {
		Response
		OEMInfo OEMInfo
	}

	err := c.CmdWithResult(ctx, CmdSystemInfoRequest, "OEMInfo", &result)
	if err != nil {
		return nil, err
	}

	return &result.OEMInfo, err
}

func (c *Client) StorageInfo(ctx context.Context) ([]StorageInfo, error) {
	var result struct {
		Response
		StorageInfo []StorageInfo
	}

	err := c.CmdWithResult(ctx, CmdSystemInfoRequest, "StorageInfo", &result)
	if err != nil {
		return nil, err
	}

	return result.StorageInfo, err
}

func (c *Client) WorkState(ctx context.Context) (*WorkState, error) {
	var result struct {
		Response
		WorkState WorkState
	}

	err := c.CmdWithResult(ctx, CmdSystemInfoRequest, "WorkState", &result)
	if err != nil {
		return nil, err
	}

	return &result.WorkState, err
}
