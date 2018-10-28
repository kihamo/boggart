package mikrotik

import (
	"context"
)

type InterfaceStats struct {
	Id         string `json:".id"`
	Name       string `json:"name"`
	MacAddress string `json:"mac-address,omitempty"`
	TXByte     uint64 `json:"tx-byte"`
	RXByte     uint64 `json:"rx-byte"`
}

type InterfaceWirelessRegistrationTable struct {
	Id         string `json:".id"`
	MacAddress string `json:"mac-address"`
	Interface  string `json:"interface"`
	Bytes      string `json:"bytes"`
}

func (c *Client) InterfaceStats(ctx context.Context) (result []InterfaceStats, err error) {
	err = c.doConvert(ctx, []string{"/interface/print", "stats"}, &result)
	return result, err
}

func (c *Client) InterfaceWirelessRegistrationTable(ctx context.Context) (result []InterfaceWirelessRegistrationTable, err error) {
	err = c.doConvert(ctx, []string{"/interface/wireless/registration-table/print"}, &result)
	return result, err
}
