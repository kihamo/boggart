package boggart // import "github.com/kihamo/boggart"

//go:generate /bin/bash -c "find components/boggart/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "find components/mqtt/internal/locales/ -name \\*.po -execdir /bin/bash -c 'msgfmt {} -o `basename {} .po`.mo' '{}' \\;"
//go:generate /bin/bash -c "goimports -w `find . -type f -name '*.go' -not -path './vendor/*' -not -name 'bindata_assetfs.go'`"
//go:generate /bin/bash -c "cd components/boggart/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal templates/... assets/... locales/..."
//go:generate /bin/bash -c "cd components/boggart && enumer -type=BindStatus -trimprefix=BindStatus -output=bind_status_enumer.go"
//go:generate /bin/bash -c "cd components/mqtt/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal templates/... locales/..."
//go:generate /bin/bash -c "cd components/voice/internal && go-bindata-assetfs -ignore='(.*?[.]po$)' -o ./bindata_assetfs.go -pkg=internal templates/..."
//go:generate /bin/bash -c "cd components/voice/players && enumer -type=Status -trimprefix=Status -output=status_enumer.go -transform=snake"
//go:generate /bin/bash -c "easyjson components/boggart/internal/handlers/manager.go"
