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

func (t Topic) Replace(replaces map[string]string) Topic {
	oldnew := make([]string, 0, len(replaces)*2)

	for k, v := range replaces {
		oldnew = append(oldnew, k, v)
	}

	topic := strings.NewReplacer(oldnew...).Replace(t.String())
	return Topic(topic)
}

func (t Topic) Split() (result []string) {
	route := strings.TrimRight(t.String(), "/")
	if strings.HasPrefix(route, "$share") {
		result = strings.Split(route, "/")[2:]
	} else {
		result = strings.Split(route, "/")
	}

	return result
}

func NameReplace(name string) string {
	name = strings.ToLower(name)
	name = strings.Trim(name, TopicSeparator)
	name = strings.Join(strings.Fields(name), " ")

	return replacerMQTTName.Replace(name)
}
