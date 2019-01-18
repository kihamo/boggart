package voice

import (
	"strings"
)

var messageReplacer = strings.NewReplacer(
	":", "",
	"[", "",
	"]", "",
	"-", "",
	";", "",
	")", "",
	"(", "",
	"\"", "",
	"'", "",
)

func TrimMessage(message string) string {
	message = strings.TrimSpace(message)
	message = strings.ToLower(message)
	message = strings.TrimFunc(message, func(r rune) bool {
		switch r {
		case '.', '!', '?', ',', ';', ':', '-', '(', ')', '"', '\'', '[', ']':
			return true
		}

		return false
	})
	message = messageReplacer.Replace(message)
	message = strings.Join(strings.Fields(message), " ")

	return message
}
