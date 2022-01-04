package mosenergosbyt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/kihamo/boggart/providers/integratorit/internal"
)

var (
	ErrProviderMethodNotSupported = errors.New("command not supported method")
)

type HasSupportCurrentBalance interface {
	CurrentBalance(ctx context.Context) (float64, error)
}

type HasSupportBills interface {
	Bills(ctx context.Context, dateStart, dateEnd time.Time) ([]Bill, error)
	BillDownload(ctx context.Context, bill Bill, writer io.Writer) error
}

type Provider interface {
}

func NewProvider(client *internal.Client, account *Account) (Provider, error) {
	switch account.ProviderID {
	case ProviderIDMosEnergoSbyt:
		return NewProvideMosEnergoSbyt(client, account), nil
	case ProviderIDMosOblERC:
		return NewProvideMosOblERC(client, account), nil
	}

	return nil, fmt.Errorf("unknown provider ID %d", account.ProviderID)
}
