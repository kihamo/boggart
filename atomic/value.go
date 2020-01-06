package atomic

import (
	"sync/atomic"
)

type Value struct {
	atomic.Value
}
