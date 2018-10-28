package mikrotik

import (
	"context"
)

type IPARP struct {
	Id         string `json:".id"`
	Address    string `json:"address"`
	MacAddress string `json:"mac-address"`
	Interface  string `json:"interface"`
	Published  bool   `json:"published"`
	Invalid    bool   `json:"invalid"`
	DHCP       bool   `json:"DHCP"`
	Dynamic    bool   `json:"dynamic"`
	Complete   bool   `json:"complete"`
	Disabled   bool   `json:"disabled"`
	Comment    string `json:"comment,omitempty"`
}

type IPDHCPServerLease struct {
	Id           string `json:".id"`
	Address      string `json:"address"`
	MacAddress   string `json:"mac-address"`
	AddressLists string `json:"address-lists"`
	Server       string `json:"server"`
	DHCPOption   string `json:"dhcp-option"`
	Status       string `json:"status"`
	LastSeen     string `json:"last-seen"`
	Radius       bool   `json:"radius"`
	Dynamic      bool   `json:"dynamic"`
	Blocked      bool   `json:"blocked"`
	Disabled     bool   `json:"disabled"`
	Comment      string `json:"comment,omitempty"`
}

type IPDNSStatic struct {
	Id       string `json:".id"`
	Name     string `json:"name"`
	Regexp   string `json:"regexp"`
	Address  string `json:"address"`
	TTL      string `json:"ttl"`
	Dynamic  bool   `json:"dynamic"`
	Disabled bool   `json:"disabled"`
	Comment  string `json:"comment,omitempty"`
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
