package openhab

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
)

type Parameters map[string]interface{}

func newParameters() *Parameters {
	return &Parameters{}
}

func (p *Parameters) Set(key string, value interface{}) {
	if p != nil {
		(*p)[key] = value
	}
}

func (p *Parameters) Delete(key string) {
	if p != nil {
		delete(*p, key)
	}
}

func (p *Parameters) String() string {
	if p == nil || len(*p) == 0 {
		return ""
	}

	parameters := make([]string, 0, len(*p))

	for key, value := range *p {
		switch v := value.(type) {
		case string:
			value = "\"" + v + "\""
		case mqtt.Topic:
			value = "\"" + v.String() + "\""
		case bool:
			if v {
				value = "true"
			} else {
				value = "false"
			}
		}

		parameters = append(parameters, fmt.Sprintf("%s=%v", key, value))
	}

	sort.Strings(parameters)

	return strings.Join(parameters, ",")
}
