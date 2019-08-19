package xmeye

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

func (c *Client) SystemInfo(ctx context.Context) (*internal.SystemInfo, error) {
	var result struct {
		internal.Response
		SystemInfo internal.SystemInfo
	}

	err := c.CmdWithResult(ctx, CmdSystemInfoRequest, "SystemInfo", &result)
	if err != nil {
		return nil, err
	}

	return &result.SystemInfo, err
}

func (c *Client) OEMInfo(ctx context.Context) (*internal.OEMInfo, error) {
	var result struct {
		internal.Response
		OEMInfo internal.OEMInfo
	}

	err := c.CmdWithResult(ctx, CmdSystemInfoRequest, "OEMInfo", &result)
	if err != nil {
		return nil, err
	}

	return &result.OEMInfo, err
}

func (c *Client) StorageInfo(ctx context.Context) (*internal.StorageInfo, error) {
	var result struct {
		internal.Response
		StorageInfo internal.StorageInfo
	}

	err := c.CmdWithResult(ctx, CmdSystemInfoRequest, "StorageInfo", &result)
	if err != nil {
		return nil, err
	}

	return &result.StorageInfo, err
}

func (c *Client) WorkState(ctx context.Context) (*internal.WorkState, error) {
	var result struct {
		internal.Response
		WorkState internal.WorkState
	}

	err := c.CmdWithResult(ctx, CmdSystemInfoRequest, "WorkState", &result)
	if err != nil {
		return nil, err
	}

	return &result.WorkState, err
}
