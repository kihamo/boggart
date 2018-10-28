package mikrotik

import (
	"context"
	"errors"
)

type SystemHealth struct {
	Voltage     float64 `json:"voltage"`
	Temperature uint64  `json:"temperature"`
}

type SystemRouterboard struct {
	Routerboard     bool   `json:"routerboard"`
	Model           string `json:"model"`
	SerialNumber    string `json:"serial-number"`
	FirmwareType    string `json:"firmware-type"`
	FactoryFirmware string `json:"factory-firmware"`
	CurrentFirmware string `json:"current-firmware"`
	UpgradeFirmware string `json:"upgrade-firmware"`
}

type SystemResource struct {
	Uptime               string  `json:"uptime"`
	Version              string  `json:"version"`
	BuildTime            string  `json:"build-time"`
	FactorySoftware      string  `json:"factory-software"`
	FreeMemory           uint64  `json:"free-memory"`
	TotalMemory          uint64  `json:"total-memory"`
	CPU                  string  `json:"cpu"`
	CPUCount             uint64  `json:"cpu-count"`
	CPUFrequency         uint64  `json:"cpu-frequency"`
	CPULoad              uint64  `json:"cpu-load"`
	FreeHDDSpace         uint64  `json:"free-hdd-space"`
	TotalHDDSpace        uint64  `json:"total-hdd-space"`
	WriteSectSinceReboot float64 `json:"write-sect-since-reboot"`
	WriteSectTotal       float64 `json:"write-sect-total"`
	BadBlocks            float64 `json:"bad-blocks"`
	ArchitectureName     string  `json:"architecture-name"`
	BoardName            string  `json:"board-name"`
	Platform             string  `json:"platform"`
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

func (c *Client) SystemHealth(ctx context.Context) (*SystemHealth, error) {
	var result []*SystemHealth

	err := c.doConvert(ctx, []string{"/system/health/print"}, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("Empty response")
	}

	return result[0], nil
}

func (c *Client) SystemRouterboard(ctx context.Context) (*SystemRouterboard, error) {
	var result []*SystemRouterboard

	err := c.doConvert(ctx, []string{"/system/routerboard/print"}, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("Empty response")
	}

	return result[0], nil
}

func (c *Client) SystemResource(ctx context.Context) (*SystemResource, error) {
	var result []*SystemResource

	err := c.doConvert(ctx, []string{"/system/resource/print"}, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("Empty response")
	}

	return result[0], nil
}

func (c *Client) SystemDisk(ctx context.Context) (result []*SystemDisk, err error) {
	err = c.doConvert(ctx, []string{"/disk/print"}, &result)
	return result, err
}
