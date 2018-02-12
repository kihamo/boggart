package boggart // import "github.com/kihamo/boggart"

//go:generate /bin/bash -c "goimports -w `find . -type f -name '*.go' -not -path './vendor/*' -not -name 'bindata_assetfs.go'`"
//go:generate /bin/bash -c "cd components/boggart/internal && go-bindata-assetfs -pkg=internal templates/..."
//go:generate /bin/bash -c "cd components/boggart && enumer -type=DeviceType -trimprefix=DeviceType -output=device_type_enumer.go"
