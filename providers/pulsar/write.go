package pulsar

import (
	"time"
)

func (d *HeatMeter) WriteDatetime() (err error) {
	t := time.Now().In(d.options.location)

	request := NewPacket().
		WithFunction(FunctionWriteTime).
		WithPayload(TimeToBytes(t))

	_, err = d.Invoke(request)

	return err
}
