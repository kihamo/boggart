```
cd $GOPATH/src/github.com/kihamo/boggart/cmd/boggart && go build && rm -rf $GOPATH/bin/boggart && go install github.com/kihamo/boggart/cmd/boggart
sudo systemctl restart boggart.service && tail -f /var/log/syslog | grep 'boggart\['
```