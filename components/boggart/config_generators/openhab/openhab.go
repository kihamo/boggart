package openhab

import (
	"strings"
)

const (
	BindingID = "mqtt"
)

var replacerID = strings.NewReplacer(
	":", "_",
	"{", "_",
	"}", "_",
	"[", "_",
	"]", "_",
	"@", "_",
	" ", "_",
	"\"", "_",
	"-", "_",
)

func IDReplace(id string) string {
	return replacerID.Replace(id)
}
