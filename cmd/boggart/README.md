```
cp /home/pi/go/src/github.com/kihamo/boggart/cmd/boggart/boggart.service /lib/systemd/system/boggart.service
sudo systemctl daemon-reload
sudo journalctl -f -u boggart.service
```

```
cd $GOPATH/src/github.com/kihamo/boggart/cmd/boggart && go build && rm -rf $GOPATH/bin/boggart && go install github.com/kihamo/boggart/cmd/boggart
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
```