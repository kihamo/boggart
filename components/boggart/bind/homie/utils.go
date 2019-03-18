package homie

import (
	"sync"
)

func syncMapToMap(m *sync.Map) map[string]interface{} {
	result := make(map[string]interface{})

	m.Range(func(key, value interface{}) bool {
		result[key.(string)] = value
		return true
	})

	return result
}
