package xmeye

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

func (c *Client) DiskManager(partNumber uint64, action diskManagerAction) error {
	_, err := c.Call(CmdDiskManagerRequest, map[string]interface{}{
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

func (c *Client) DiskManagerPartition(partNumber, record, snapshot uint64) error {
	_, err := c.Call(CmdDiskManagerRequest, map[string]interface{}{
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

func (c *Client) DiskManagerSetType(partNumber uint64, typ diskManagerType) error {
	_, err := c.Call(CmdDiskManagerRequest, map[string]interface{}{
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
