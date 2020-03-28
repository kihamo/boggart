package mikrotik

import (
	"context"

	"github.com/kihamo/boggart/types"
)

type IPARP struct {
	ID         string             `mapstructure:".id"`
	Address    types.IP           `mapstructure:"address"`
	MacAddress types.HardwareAddr `mapstructure:"mac-address,omitempty"`
	Interface  string             `mapstructure:"interface"`
	Published  bool               `mapstructure:"published"`
	Invalid    bool               `mapstructure:"invalid"`
	DHCP       bool               `mapstructure:"DHCP"`
	Dynamic    bool               `mapstructure:"dynamic"`
	Complete   bool               `mapstructure:"complete"`
	Disabled   bool               `mapstructure:"disabled"`
	Comment    string             `mapstructure:"comment,omitempty"`
}

type IPDHCPServerLease struct {
	ID           string             `mapstructure:".id"`
	Address      types.IP           `mapstructure:"address"`
	MacAddress   types.HardwareAddr `mapstructure:"mac-address"`
	AddressLists string             `mapstructure:"address-lists"`
	Server       string             `mapstructure:"server"`
	DHCPOption   string             `mapstructure:"dhcp-option"`
	Status       string             `mapstructure:"status"`
	LastSeen     string             `mapstructure:"last-seen"`
	Radius       bool               `mapstructure:"radius"`
	Dynamic      bool               `mapstructure:"dynamic"`
	Blocked      bool               `mapstructure:"blocked"`
	Disabled     bool               `mapstructure:"disabled"`
	Comment      string             `mapstructure:"comment,omitempty"`
}

type IPDNSStatic struct {
	ID       string   `mapstructure:".id"`
	Name     string   `mapstructure:"name"`
	Regexp   string   `mapstructure:"regexp"`
	Address  types.IP `mapstructure:"address"`
	TTL      string   `mapstructure:"ttl"`
	Dynamic  bool     `mapstructure:"dynamic"`
	Disabled bool     `mapstructure:"disabled"`
	Comment  string   `mapstructure:"comment,omitempty"`
}

func (c *Client) IPARP(ctx context.Context) (result []IPARP, err error) {
	err = c.doConvert(ctx, []string{"/ip/arp/print"}, &result)
	return result, err
}

func (c *Client) IPDHCPServerLease(ctx context.Context) (result []IPDHCPServerLease, err error) {
	err = c.doConvert(ctx, []string{"/ip/dhcp-server/lease/print"}, &result)
	return result, err
}

func (c *Client) IPDNSStatic(ctx context.Context) (result []IPDNSStatic, err error) {
	err = c.doConvert(ctx, []string{"/ip/dns/static/print"}, &result)
	return result, err
}
