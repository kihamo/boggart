package rkcm

import (
	"net"
	"strconv"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/components/boggart/providers/rkcm/client"
)

// http://uk-gkh.org/muprcmytishi/GetMobileMenu.ashx?appVersion=2.112
// http://uk-gkh.org/muprcmytishi/MobileAPI/AuthenticateAccount.ashx?phone={phone}&pwd={password}
// http://uk-gkh.org/muprcmytishi/RegisterClientDevice.ashx?cid=13076&did={token}&os=Android&version=28&model={model}&isMobAcc=1
// http://uk-gkh.org/muprcmytishi/DataExport.ashx?table=Support_RequestTypes
// http://uk-gkh.org/muprcmytishi/GetPayments.ashx?login={phone}&pwd={password}
// http://uk-gkh.org/muprcmytishi/CheckQuestionsNeedUpdate.ashx?count=0
// http://uk-gkh.org/muprcmytishi/MobileAPI/GetPays.ashx?phone={phone}
// http://uk-gkh.org/muprcmytishi/GetHousesWebCams.ashx?ident={ident}
// http://uk-gkh.org/muprcmytishi/GetAdditionalServices.ashx?login={phone}&pwd={password}

const (
	defaultHost = "uk-gkh.org"
	defaultPort = 80
)

type Client struct {
	*client.RKCM

	phone    string
	password string
}

func New(phone, password string, debug bool, logger logger.Logger) *Client {
	cfg := client.DefaultTransportConfig().WithHost(net.JoinHostPort(defaultHost, strconv.FormatInt(defaultPort, 10)))
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		// что бы скачивались файлы с изображениями
		rt.Consumers["image/jpeg"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/png"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/gif"] = runtime.ByteStreamConsumer()

		rt.Consumers["application/octet-stream"] = runtime.XMLConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return &Client{
		RKCM: cl,

		phone:    phone,
		password: password,
	}
}
