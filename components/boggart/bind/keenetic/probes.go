package keenetic

import (
	"context"

	"github.com/kihamo/boggart/providers/keenetic/client/user"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	_, err = b.client.User.GetAuth(user.NewGetAuthParamsWithContext(ctx))

	if _, ok := err.(*user.GetAuthUnauthorized); ok {
		if u := b.config().Address.User; u != nil {
			password, _ := u.Password()

			return b.client.Login(ctx, u.Username(), password)
		}
	}

	return err
}
