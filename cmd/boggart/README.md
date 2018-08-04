```
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/boggart/boggart.service /lib/systemd/system/boggart.service
sudo systemctl daemon-reload
sudo journalctl -f -u boggart.service
```

```
GOARM=7 gox -output="cmd/boggart/boggart" -osarch="linux/arm" ./cmd/boggart/
cp -f $GOPATH/src/github.com/kihamo/boggart/cmd/boggart/boggart $GOPATH/bin/boggart

cd $GOPATH/src/github.com/kihamo/boggart/cmd/boggart && go build -v && rm -rf $GOPATH/bin/boggart && go install github.com/kihamo/boggart/cmd/boggart

sudo systemctl stop boggart.service
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service

sudo echo "22" > /sys/class/gpio/unexport
sudo echo "27" > /sys/class/gpio/unexport
```