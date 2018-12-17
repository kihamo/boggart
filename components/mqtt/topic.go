package mqtt

import (
	"fmt"
	"strings"
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
				parts[i] = fmt.Sprintf("%v", arg)
				break
			}
		}
	}

	return strings.Join(parts, "/")
}
