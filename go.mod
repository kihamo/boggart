module github.com/kihamo/boggart

go 1.13

require (
	cloud.google.com/go/firestore v1.2.0 // indirect
	firebase.google.com/go v3.12.0+incompatible
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496
	github.com/barnybug/go-cast v0.0.0-20190910160619-d2aa97f56d4e
	github.com/bieber/barcode v0.0.0-20190908000948-a94135955bb1
	github.com/denisbrodbeck/machineid v1.0.1
	github.com/eclipse/paho.mqtt.golang v1.2.1-0.20200511074540-4c98a2381d16
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/faiface/beep v1.0.2
	github.com/fsnotify/fsnotify v1.4.9
	github.com/ghthor/gowol v0.0.0-20180205141434-eb42ead1b24e
	github.com/go-ble/ble v0.0.0-20200120171844-0a73a9da88eb
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/runtime v0.19.14
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.8
	github.com/go-openapi/validate v0.19.7
	github.com/goburrow/serial v0.1.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.5
	github.com/gorilla/websocket v1.4.2
	github.com/hajimehoshi/oto v0.5.4
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/influxdata/influxdb v1.8.0
	github.com/kihamo/go-workers v2.1.7+incompatible
	github.com/kihamo/shadow v0.0.0-20200412185102-25db29c8fa3b
	github.com/kihamo/snitch v0.0.0-20200412182537-3478a87783e1
	github.com/llgcode/draw2d v0.0.0-20200110163050-b96d8208fcfc
	github.com/mailru/easyjson v0.7.1
	github.com/mitchellh/mapstructure v1.2.2
	github.com/mmcloughlin/geohash v0.9.0
	github.com/mourner/suncalc-go v0.0.0-20141021103505-77cea98fd55e
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/snabb/webostv v0.0.1
	github.com/sparrc/go-ping v0.0.0-20190613174326-4e5b6552494c
	github.com/yryz/ds18b20 v0.0.0-20180211073435-3cf383a40624
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.14.1
	gocv.io/x/gocv v0.22.0 // indirect
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/api v0.20.0
	google.golang.org/grpc v1.28.0
	gopkg.in/mcuadros/go-syslog.v2 v2.3.0
	gopkg.in/routeros.v2 v2.0.0-20190905230420-1bbf141cdd91
	gopkg.in/telegram-bot-api.v4 v4.6.4
	gopkg.in/yaml.v2 v2.2.8
	periph.io/x/periph v3.6.2+incompatible
)

replace github.com/barnybug/go-cast v0.0.0-20190910160619-d2aa97f56d4e => github.com/kihamo/go-cast v0.0.0-20190130214031-2bd907ad55c2

replace github.com/sparrc/go-ping => github.com/kihamo/go-ping v0.0.0-20200405124135-bc7921838e0d

replace github.com/kihamo/shadow v0.0.0-20200412160130-a95972d7c957 => /Users/kihamo/go/src/github.com/kihamo/shadow
