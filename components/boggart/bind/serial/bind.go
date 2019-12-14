package serial

import (
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/serial_network"
)

type Bind struct {
	status uint32 // 0 - default, 1 - started, 2 - failed

	boggart.BindBase

	server serial_network.Server
}

func (b *Bind) Run() error {
	go func() {
		atomic.StoreUint32(&b.status, 1)
		defer atomic.StoreUint32(&b.status, 2)

		if err := b.server.ListenAndServe(); err != nil {
			b.Logger().Error("Failed serve with error " + err.Error())
		}
	}()

	return nil
}

func (b *Bind) Close() (err error) {
	return b.server.Close()
}
