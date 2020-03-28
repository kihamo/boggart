package mikrotik

import (
	"context"

	"github.com/kihamo/boggart/types"
)

type InterfaceStats struct {
	Id         string             `mapstructure:".id"`
	Name       string             `mapstructure:"name"`
	MacAddress types.HardwareAddr `mapstructure:"mac-address,omitempty"`
	TXByte     uint64             `mapstructure:"tx-byte"`
	RXByte     uint64             `mapstructure:"rx-byte"`
}

type InterfaceWirelessRegistrationTable struct {
	Id         string             `mapstructure:".id"`
	MacAddress types.HardwareAddr `mapstructure:"mac-address"`
	Interface  string             `mapstructure:"interface"`
	Bytes      string             `mapstructure:"bytes"`
}

func (c *Client) InterfaceStats(ctx context.Context) (result []InterfaceStats, err error) {
	err = c.doConvert(ctx, []string{"/interface/print", "stats"}, &result)
	return result, err
}

func (c *Client) InterfaceWirelessRegistrationTable(ctx context.Context) (result []InterfaceWirelessRegistrationTable, err error) {
	err = c.doConvert(ctx, []string{"/interface/wireless/registration-table/print"}, &result)
	return result, err
}
