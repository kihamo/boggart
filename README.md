```
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/boggart/boggart.service /lib/systemd/system/boggart.service
sudo systemctl daemon-reload
sudo journalctl -f -u boggart.service
```

## Server
#### First
```
gox -output="cmd/server/boggart" -osarch="linux/amd64" -ldflags="-X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/server/
sudo cp -f /home/kihamo/go/src/github.com/kihamo/boggart/cmd/server/boggart.service /lib/systemd/system/boggart.service
sudo cp -f /home/kihamo/go/src/github.com/kihamo/boggart/cmd/server/boggart /usr/local/bin/boggart-server
sudo chmod +x /usr/local/bin/boggart-server
sudo systemctl daemon-reload
sudo systemctl enable boggart.service
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service
```
#### Update
```
gox -output="cmd/server/boggart" -osarch="linux/amd64" -ldflags="-X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/server/
sudo cp -f /home/kihamo/go/src/github.com/kihamo/boggart/cmd/server/boggart /usr/local/bin/boggart-server
sudo chmod +x /usr/local/bin/boggart-server
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
```

## Agent
#### First
```
GOARM=7 gox -output="cmd/agent/boggart" -osarch="linux/arm" -ldflags="-X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/agent/
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/agent/boggart.service /lib/systemd/system/boggart.service
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/agent/boggart /usr/local/bin/boggart-agent
sudo chmod +x /usr/local/bin/boggart-agent
sudo systemctl daemon-reload
sudo systemctl enable boggart.service
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service
```
#### Update
```
GOARM=7 gox -output="cmd/agent/boggart" -osarch="linux/arm"  -ldflags="-X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/agent/
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/agent/boggart /usr/local/bin/boggart-agent
sudo chmod +x /usr/local/bin/boggart-agent
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
```

## Agent Roborock
Кросс компиляция не возможна из-за ошибок, поэтому собирать надо на реальном устройстве. Подойдет Raspberry PI, так как платформа на нем аналогичная

Для библиотеки https://github.com/hajimehoshi/oto
```
apt install libasound2-dev
```
Иначе будет ошибка
```
# github.com/kihamo/boggart/vendor/github.com/hajimehoshi/oto
vendor/github.com/hajimehoshi/oto/player_linux.go:23:28: fatal error: alsa/asoundlib.h: No such file or directory
 #include <alsa/asoundlib.h>
                            ^
compilation terminated.
```

// TODO: агент должен стартовать после инициализации wifi, иначе MQTT соединение пустое

#### Build on RPi
```
cd $GOPATH/src/github.com/kihamo/boggart/cmd/roborock/
go build -ldflags "-s -w -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" -o boggart ./
```

#### First
```
sudo cp -f /home/cleaner/boggart.env /etc/default/boggart-roborock
sudo cp -f /home/cleaner/boggart.service /etc/init.d/boggart-roborock && sudo chmod +x /etc/init.d/boggart-roborock
sudo cp -f /home/cleaner/boggart /usr/local/bin/boggart-roborock && sudo chmod +x /usr/local/bin/boggart-roborock
sudo update-rc.d boggart-roborock defaults 90
sudo update-rc.d boggart-roborock enable
sudo service boggart-roborock restart
```
#### Update
```
sudo cp -f /home/cleaner/boggart.env /etc/default/boggart-roborock
sudo cp -f /home/cleaner/boggart /usr/local/bin/boggart-roborock && sudo chmod +x /usr/local/bin/boggart-roborock
sudo service boggart-roborock stop
sudo service boggart-roborock start
```

## Other

```
sudo systemctl stop boggart.service
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service

sudo echo "22" > /sys/class/gpio/unexport
sudo echo "27" > /sys/class/gpio/unexport
sudo echo "5" > /sys/class/gpio/unexport
```