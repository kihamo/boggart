package mqtt

import (
	"fmt"
	"strings"
)

const (
	TopicSeparator = "-"
)

var replacerMQTTName = strings.NewReplacer(
	":", TopicSeparator,
	"/", TopicSeparator,
	"_", TopicSeparator,
	",", TopicSeparator,
	".", TopicSeparator,
	" ", TopicSeparator,
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
	name = strings.Trim(name, TopicSeparator)
	name = strings.Join(strings.Fields(name), " ")

	return replacerMQTTName.Replace(name)
}
