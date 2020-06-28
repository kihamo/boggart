package z_stack

import (
	"context"
)

func (b *Bind) syncPermitJoin() {
	if sn := b.Meta().SerialNumber(); sn != "" {
		b.MQTT().PublishAsync(context.TODO(), b.config.TopicPermitJoinState.Format(sn), b.client.PermitJoinEnabled())
	}
}
