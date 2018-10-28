package mikrotik

import (
	"context"
)

type SystemHealth struct {
	Voltage     float64
	Temperature uint64
}

type SystemDisk struct {
	Id    string `json:".id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"type"`
	Disk  string `json:"disk"`
	Free  uint64 `json:"free"`
	Size  uint64 `json:"size"`
}

func (c *Client) SystemDisk(ctx context.Context) (result []SystemDisk, err error) {
	err = c.doConvert(ctx, []string{"/disk/print"}, &result)
	return result, err
}
