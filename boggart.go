package boggart // import "github.com/kihamo/boggart"

//go:generate /bin/bash -c "find components/boggart/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/alsa/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/astro/sun/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/broadlink/sp3s/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/chromecast/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/ds18b20/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/gpio/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/hikvision/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/homie/esp/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/mercury/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/lg_webos/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/nut/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/pulsar/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/boggart/bind/xiaomi/roborock/miio/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/mqtt/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/openhab/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/storage/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/syslog/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/voice/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "goimports -w `find . -type f -name '*.go' -not -path './vendor/*' -not -name 'bindata_assetfs.go'`"
//go:generate /bin/bash -c "cd components/boggart/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal -nometadata -nomemcopy templates/... assets/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/alsa && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=alsa -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/astro/sun && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=sun -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/broadlink/sp3s && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=sp3s -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/chromecast && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=chromecast -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/ds18b20 && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=ds18b20 -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/gpio && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=gpio -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/hikvision && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=hikvision -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/homie/esp && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=esp -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/lg_webos && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=lg_webos -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/mercury && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=mercury -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/nut && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=nut -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/pulsar && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=pulsar -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/xiaomi/roborock/miio && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=miio -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart && enumer -type=BindStatus -trimprefix=BindStatus -output=bind_status_enumer.go"
//go:generate /bin/bash -c "cd components/mqtt/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal -nometadata -nomemcopy templates/... locales/..."
//go:generate /bin/bash -c "cd components/openhab/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal -nometadata -nomemcopy locales/..."
//go:generate /bin/bash -c "cd components/storage/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal -nometadata -nomemcopy locales/..."
//go:generate /bin/bash -c "cd components/syslog/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal -nometadata -nomemcopy locales/..."
//go:generate /bin/bash -c "cd components/voice/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal -nometadata -nomemcopy locales/..."
//go:generate /bin/bash -c "cd components/boggart/bind/alsa && enumer -type=Status -trimprefix=Status -output=status_enumer.go -transform=snake"
//go:generate /bin/bash -c "easyjson components/boggart/internal/handlers/manager.go"
