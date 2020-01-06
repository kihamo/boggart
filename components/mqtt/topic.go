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

func (t Topic) Format(args ...interface{}) Topic {
	parts := t.Split()

	for _, arg := range args {
		for i, topic := range parts {
			if topic == "+" {
				parts[i] = NameReplace(fmt.Sprintf("%v", arg))
				break
			}
		}
	}

	topic := strings.Join(parts, "/")
	return Topic(topic)
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

func (t Topic) ValidAsPublishTopic() bool {
	s := t.String()

	if strings.Contains(s, "+") {
		return false
	}

	if strings.Contains(s, "#") {
		return false
	}

	return true
}

func (t Topic) ValidAsSubscribeTopic() bool {
	if s := t.String(); strings.Contains(s, "#") {
		if s[len(s)-1:] != "#" {
			return false
		}
	}

	return true
}

func (t Topic) match(routes1 []string, routes2 []string) bool {
	if len(routes1) == 0 {
		return len(routes2) == 0
	}

	if len(routes2) == 0 {
		return routes1[0] == "#"
	}

	if routes1[0] == "#" {
		return true
	}

	if (routes1[0] == "+") || (routes1[0] == routes2[0]) {
		return t.match(routes1[1:], routes2[1:])
	}

	return false
}

func (t Topic) IsInclude(topic Topic) bool {
	if t.String() == topic.String() {
		return true
	}

	topic = Topic(strings.TrimRight(topic.String(), "/"))
	return t.match(t.Split(), topic.Split())
}

func NameReplace(name string) string {
	name = strings.ToLower(name)
	name = strings.Trim(name, TopicSeparator)
	name = strings.Join(strings.Fields(name), " ")

	return replacerMQTTName.Replace(name)
}
