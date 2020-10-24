package atomic

import (
	"sync/atomic"
)

type Value struct {
	atomic.Value
}

func NewValue() *Value {
	return &Value{}
}
