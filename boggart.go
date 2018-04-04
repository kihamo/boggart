package boggart // import "github.com/kihamo/boggart"

//go:generate /bin/bash -c "find components/boggart/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "goimports -w `find . -type f -name '*.go' -not -path './vendor/*' -not -name 'bindata_assetfs.go'`"
//go:generate /bin/bash -c "cd components/boggart/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -pkg=internal templates/... locales/..."
//go:generate /bin/bash -c "cd components/boggart && enumer -type=DeviceType -trimprefix=DeviceType -output=device_type_enumer.go"
//go:generate /bin/bash -c "cd components/boggart && enumer -type=DeviceId -trimprefix=DeviceId -output=device_id_enumer.go -transform=snake"
