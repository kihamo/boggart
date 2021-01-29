package mikrotik

import (
	"context"

	"github.com/kihamo/boggart/types"
)

type Interface struct {
	MacAddress       types.HardwareAddr `mapstructure:"mac-address"`
	ID               string             `mapstructure:".id"`
	DefaultName      string             `mapstructure:"default-name"`
	LastLinkDownTime string             `mapstructure:"last-link-down-time"`
	LastLinkUpTime   string             `mapstructure:"last-link-up-time"`
	Name             string             `mapstructure:"name"`
	Type             string             `mapstructure:"type"`
	ActualMTU        uint64             `mapstructure:"actual-mtu"`
	FpRxByte         uint64             `mapstructure:"fp-rx-byte"`
	FpRxPacket       uint64             `mapstructure:"fp-rx-packet"`
	FpTxByte         uint64             `mapstructure:"fp-tx-byte"`
	FpTxPacket       uint64             `mapstructure:"fp-tx-packet"`
	L2MTU            uint64             `mapstructure:"l2mtu"`
	LinkDowns        uint64             `mapstructure:"link-downs"`
	MaxL2MTU         uint64             `mapstructure:"max-l2mtu"`
	RxByte           uint64             `mapstructure:"rx-byte"`
	RxDrop           uint64             `mapstructure:"rx-drop"`
	RxError          uint64             `mapstructure:"rx-error"`
	RxPacket         uint64             `mapstructure:"rx-packet"`
	TxByte           uint64             `mapstructure:"tx-byte"`
	TxDrop           uint64             `mapstructure:"tx-drop"`
	TxError          uint64             `mapstructure:"tx-error"`
	TxPacket         uint64             `mapstructure:"tx-packet"`
	TxQueueDrop      uint64             `mapstructure:"queue-drop"`
	Disabled         bool               `mapstructure:"disabled"`
	Running          bool               `mapstructure:"running"`
	Slave            bool               `mapstructure:"slave"`
}

type InterfaceL2TP struct {
	ID       string `mapstructure:".id"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Running  bool   `mapstructure:"running"`
	Disabled bool   `mapstructure:"disabled"`
}

type InterfaceStats struct {
	ID         string             `mapstructure:".id"`
	Name       string             `mapstructure:"name"`
	MacAddress types.HardwareAddr `mapstructure:"mac-address,omitempty"`
	TXByte     uint64             `mapstructure:"tx-byte"`
	RXByte     uint64             `mapstructure:"rx-byte"`
}

type InterfaceWirelessRegistrationTable struct {
	ID         string             `mapstructure:".id"`
	MacAddress types.HardwareAddr `mapstructure:"mac-address"`
	Interface  string             `mapstructure:"interface"`
	Bytes      string             `mapstructure:"bytes"`
}

func (c *Client) Interfaces(ctx context.Context) (result []Interface, err error) {
	err = c.doConvert(ctx, []string{"/interface/print"}, &result)
	return result, err
}

func (c *Client) InterfaceStats(ctx context.Context) (result []InterfaceStats, err error) {
	err = c.doConvert(ctx, []string{"/interface/print", "stats"}, &result)
	return result, err
}

func (c *Client) InterfaceL2TPServer(ctx context.Context) (result []InterfaceL2TP, err error) {
	err = c.doConvert(ctx, []string{"/interface/l2tp-server/print"}, &result)
	return result, err
}

func (c *Client) InterfaceWirelessRegistrationTable(ctx context.Context) (result []InterfaceWirelessRegistrationTable, err error) {
	err = c.doConvert(ctx, []string{"/interface/wireless/registration-table/print"}, &result)
	return result, err
}
