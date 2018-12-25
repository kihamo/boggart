package mqtt

import (
	"fmt"
	"strings"
)

var replacerMQTTName = strings.NewReplacer(
	":", "-",
	"/", "-",
	"_", "-",
	",", "-",
	".", "-",
)

type Topic string

func (t Topic) String() string {
	return string(t)
}

func (t Topic) Format(args ...interface{}) string {
	parts := RouteSplit(t.String())

	for _, arg := range args {
		for i, topic := range parts {
			if topic == "+" {
				parts[i] = NameReplace(fmt.Sprintf("%v", arg))
				break
			}
		}
	}

	return strings.Join(parts, "/")
}

func NameReplace(name string) string {
	name = strings.ToLower(name)
	return replacerMQTTName.Replace(name)
}
