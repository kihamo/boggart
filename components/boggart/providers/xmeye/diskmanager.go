package xmeye

import (
	"context"
)

type diskManagerAction string
type diskManagerType string

const (
	DiskManagerActionClear     diskManagerAction = "Clear"
	DiskManagerActionPartition diskManagerAction = "Partition"
	DiskManagerActionRecover   diskManagerAction = "Recover"
	DiskManagerActionSetType   diskManagerAction = "SetType"

	DiskManagerTypeReadOnly  diskManagerType = "ReadOnly"
	DiskManagerTypeReadWrite diskManagerType = "ReadWrite"
)

func (c *Client) DiskManager(ctx context.Context, partNumber uint64, action diskManagerAction) error {
	_, err := c.Call(ctx, CmdDiskManagerRequest, map[string]interface{}{
		"Name":      "OPStorageManager",
		"SessionID": c.sessionIDAsString(),
		"OPStorageManager": map[string]interface{}{
			"Action":   string(action),
			"PartNo":   partNumber,
			"SerialNo": 0,
		},
	})

	return err
}

func (c *Client) DiskManagerPartition(ctx context.Context, partNumber, record, snapshot uint64) error {
	_, err := c.Call(ctx, CmdDiskManagerRequest, map[string]interface{}{
		"Name":      "OPStorageManager",
		"SessionID": c.sessionIDAsString(),
		"OPStorageManager": map[string]interface{}{
			"Action":   string(DiskManagerActionPartition),
			"PartNo":   partNumber,
			"SerialNo": 0,
			"PartitionSize": map[string]uint64{
				"Record":   record,
				"SnapShot": snapshot,
			},
		},
	})

	return err
}

func (c *Client) DiskManagerSetType(ctx context.Context, partNumber uint64, typ diskManagerType) error {
	_, err := c.Call(ctx, CmdDiskManagerRequest, map[string]interface{}{
		"Name":      "OPStorageManager",
		"SessionID": c.sessionIDAsString(),
		"OPStorageManager": map[string]interface{}{
			"Action":   string(DiskManagerActionSetType),
			"PartNo":   partNumber,
			"SerialNo": 0,
			"Type":     string(typ),
		},
	})

	return err
}
