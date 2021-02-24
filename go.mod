module github.com/kihamo/boggart

go 1.13

require (
	cloud.google.com/go/firestore v1.2.0 // indirect
	firebase.google.com/go v3.12.0+incompatible
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535
	github.com/barnybug/go-cast v0.0.0-20190910160619-d2aa97f56d4e
	github.com/bieber/barcode v0.0.0-20201127170204-1d90414c63eb
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/eclipse/paho.mqtt.golang v1.2.1-0.20200609161119-ca94c5368c77
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/faiface/beep v1.0.2
	github.com/fsnotify/fsnotify v1.4.9
	github.com/ghthor/gowol v0.0.0-20180205141434-eb42ead1b24e
	github.com/go-ble/ble v0.0.0-20200407180624-067514cd6e24
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/runtime v0.19.16
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/go-openapi/validate v0.19.10
	github.com/goburrow/serial v0.1.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/websocket v1.4.2
	github.com/hajimehoshi/oto v0.5.4
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/influxdata/influxdb v1.8.0
	github.com/influxdata/influxdb-client-go v1.4.0
	github.com/kihamo/shadow v0.0.0-20210224231027-75d4f243ad4d
	github.com/kihamo/snitch v0.0.0-20200412182537-3478a87783e1
	github.com/llgcode/draw2d v0.0.0-20200930101115-bfaf5d914d1e
	github.com/mailru/easyjson v0.7.6
	github.com/mitchellh/mapstructure v1.3.2
	github.com/mmcloughlin/geohash v0.9.0
	github.com/mourner/suncalc-go v0.0.0-20141021103505-77cea98fd55e
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.9.1
	github.com/snabb/webostv v0.0.1
	github.com/sparrc/go-ping v0.0.0-20190613174326-4e5b6552494c
	github.com/yryz/ds18b20 v0.0.0-20180211073435-3cf383a40624
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.14.1
	gocv.io/x/gocv v0.25.0 // indirect
	golang.org/x/net v0.0.0-20200904194848-62affa334b73
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20200610111108-226ff32320da // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e
	golang.org/x/tools v0.0.0-20200911040025-d179df38ff46 // indirect
	google.golang.org/api v0.20.0
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
	gopkg.in/routeros.v2 v2.0.0-20190905230420-1bbf141cdd91
	gopkg.in/telegram-bot-api.v4 v4.6.4
	gopkg.in/yaml.v2 v2.3.0
	periph.io/x/periph v3.6.4+incompatible
)

replace github.com/barnybug/go-cast v0.0.0-20190910160619-d2aa97f56d4e => github.com/kihamo/go-cast v0.0.0-20190130214031-2bd907ad55c2

replace github.com/sparrc/go-ping => github.com/kihamo/go-ping v0.0.0-20200405124135-bc7921838e0d
