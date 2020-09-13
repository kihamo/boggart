```
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/boggart/boggart.service /lib/systemd/system/boggart.service
sudo systemctl daemon-reload
sudo journalctl -f -u boggart.service
```

## Server
#### First
```
gox -output="cmd/server/boggart" -osarch="linux/amd64" -ldflags="-X 'main.Name=Boggart Server' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/server/
sudo cp -f /home/kihamo/go/src/github.com/kihamo/boggart/cmd/server/boggart.service /lib/systemd/system/boggart.service
sudo cp -f /home/kihamo/go/src/github.com/kihamo/boggart/cmd/server/boggart /usr/local/bin/boggart-server
sudo chmod +x /usr/local/bin/boggart-server
sudo systemctl daemon-reload
sudo systemctl enable boggart.service
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service
```
#### Update
```
gox -output="cmd/server/boggart" -osarch="linux/amd64" -ldflags="-X 'main.Name=Boggart Server' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/server/
sudo cp -f /home/kihamo/go/src/github.com/kihamo/boggart/cmd/server/boggart /usr/local/bin/boggart-server
sudo chmod +x /usr/local/bin/boggart-server
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
```

## Agent
#### First
```
GOARM=7 gox -output="cmd/agent/boggart" -osarch="linux/arm" -ldflags="-X 'main.Name=Boggart Agent' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/agent/
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/agent/boggart.service /lib/systemd/system/boggart.service
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/agent/boggart /usr/local/bin/boggart-agent
sudo chmod +x /usr/local/bin/boggart-agent
sudo systemctl daemon-reload
sudo systemctl enable boggart.service
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service
```
#### Update
###### Build for Raspberry 2
```
GOARM=7 gox -output="cmd/agent/boggart" -osarch="linux/arm" -ldflags="-X 'main.Name=Boggart Agent' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/agent/
```
###### Build for Raspberry 1
```
GOARM=6 gox -output="cmd/agent/boggart" -osarch="linux/arm" -ldflags="-X 'main.Name=Boggart Zav' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/agent/
```
###### Install
```
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/agent/boggart /usr/local/bin/boggart-agent
sudo chmod +x /usr/local/bin/boggart-agent
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
```

#### Self
```
cd $GOPATH/src/github.com/kihamo/boggart/cmd/agent/
go build -ldflags "-s -w -X 'main.Name=Boggart Villa' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" -o boggart ./
sudo cp -f $GOPATH/src/github.com/kihamo/boggart/cmd/agent/boggart /usr/local/bin/boggart-agent
sudo chmod +x /usr/local/bin/boggart-agent
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
```

#### Bluetooth
```
sudo systemctl daemon-reload
sudo systemctl enable boggart.service
sudo systemctl enable hciuart.service

sudo systemctl start bluetooth.service
sudo systemctl start hciuart.service

sudo systemctl status bluetooth.service
sudo systemctl status hciuart.service

sudo bluetoothctl

agent on
default-agent
scan on
```

В случае ошибки hci0: can't up device: connection timed out необходимо перезапустить службу
```
service bluetooth restart
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
md5sum boggart
```

#### First
```
sudo cp -f /home/cleaner/boggart.env /etc/default/boggart-roborock
sudo cp -f /home/cleaner/boggart.service /etc/init.d/boggart-roborock && sudo chmod +x /etc/init.d/boggart-roborock
sudo cp -f /home/cleaner/boggart.conf /etc/init/boggart-roborock.conf
sudo cp -f /home/cleaner/boggart /usr/local/bin/boggart-roborock && sudo chmod +x /usr/local/bin/boggart-roborock
# sudo update-rc.d boggart-roborock defaults 90
# sudo update-rc.d boggart-roborock enable
sudo service boggart-roborock restart
```
#### Update
```
md5sum /home/cleaner/boggart
sudo cp -f /home/cleaner/boggart.env /etc/default/boggart-roborock
sudo cp -f /home/cleaner/boggart /usr/local/bin/boggart-roborock && sudo chmod +x /usr/local/bin/boggart-roborock
sudo /etc/init.d/boggart-roborock stop
sudo /etc/init.d/boggart-roborock start && tail -f /var/log/boggart-roborock.log
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

```
cd /Users/kihamo/go/src/github.com/kihamo/boggart/openhab2/esphome
docker pull esphome/esphome
docker run --rm -p 6052:6052 -p 6123:6123 -v "${PWD}":/config -e ESPHOME_DASHBOARD_USE_PING=true -it esphome/esphome
```

# openhab

```
sudo systemctl stop openhab2
sudo openhab-cli clean-cache
sudo systemctl start openhab2
```

# docker

#### Swarm
Инициализация кластера
```
docker swarm init --advertise-addr 192.168.70.10
```
Что бы потом получить строку с джоинтом к кластеру
```
docker swarm join-token manager
```

Включаем не защищенные репы в маке
```
For docker desktop 2.3.x in MAC , it can be set as follows: Go to "docker" -> "preferences" -> "Docker Engine" and add the following:

insecure-registries": [
    "my private docker hub  url"
  ]
```

#### Registry
```
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

#### Сборка образа
Собираем через docker-composer
```
docker-compose build boggart-server
```

Пушим в удаленную приватную репу
```
docker-compose push boggart-server
```

Деплоем через stack (на ноде)
```
cd $GOPATH/src/github.com/kihamo/boggart/
docker stack deploy --compose-file docker-compose.yml boggart
```
Просмотр лога предпоследнего контейнера (аналог кшал лога)
```
docker service logs -f $(docker stack ps boggart -q | head -2 | tail -1)
```
```
docker service logs -f $(docker stack ps boggart -q | head -1)
```
Залезаем внутрь реплики
```
docker exec -it $(docker ps -q -f name=boggart_boggart-server) /bin/bash
```
Как правильно прокидывать устройства
https://github.com/Koenkk/zigbee2mqtt/issues/2049