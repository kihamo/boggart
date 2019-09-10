package home

//go:generate /bin/bash -c "swagger generate client -f ./swagger.yaml -t ./ --skip-validation"

import (
	"strconv"
	"strings"

	"github.com/kihamo/boggart/providers/google/home/client"
)

const (
	DefaultPort = 8008
)

func NewClient(host string) *client.GoogleHome {
	parts := strings.Split(host, ":")
	if len(parts) < 2 {
		host = strings.Join(append(parts, strconv.Itoa(DefaultPort)), ":")
	}

	cfg := client.DefaultTransportConfig().WithHost(host)
	return client.NewHTTPClientWithConfig(nil, cfg)
}
