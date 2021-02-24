```
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/boggart/boggart.service /lib/systemd/system/boggart.service
sudo systemctl daemon-reload
sudo journalctl -f -u boggart.service
```

## Server
#### First
```
gox -output="cmd/server/boggart" -osarch="linux/amd64" -ldflags="-X 'main.Name=Boggart Server' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/server/
sudo cp -f /opt/boggart/boggart.service /lib/systemd/system/boggart.service
sudo cp -f /opt/boggart/boggart /usr/local/bin/boggart-server
sudo chmod +x /usr/local/bin/boggart-server
sudo systemctl daemon-reload
sudo systemctl enable boggart.service
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service
```
#### Update
```
gox -output="cmd/server/boggart" -osarch="linux/amd64" -ldflags="-X 'main.Name=Boggart Server' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/server/
sudo cp -f /opt/boggart/boggart /usr/local/bin/boggart-server
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
git fetch --prune
git reset --hard origin/master
go build -v -ldflags "-s -w -X 'main.Name=Boggart Villa' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" -o boggart ./
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

В случае ошибки Sap driver initialization failed.
```
/etc/systemd/system/bluetooth.target.wants/bluetooth.service
ExecStart=/usr/lib/bluetooth/bluetoothd --noplugin=sap
sudo systemctl daemon-reload
sudo service bluetooth restart
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

#### zigbee2mqtt
```
cd /opt/zigbee2mqtt
npm ci
```

```
DEBUG="zigbee-herdsman*" npm start
```

#### CC2531 firmware with Rpi


##### Pins on stick CC2531
```
      /\/\/\ Antenna /\/\/\
        GND [1] | [2] V++
            ---------
Debug Clock [3] | [4] Debug Data
            ---------
  SPI Select 5  |  6  SPI Clock
            ---------
      Reset [7] |  8  SPI Data out
            ---------
   3.3V out  9  | 10  SPI Data in
            ---------
           |   USB   |
            ---------
```

##### Pins on stick RPi B
```
      ---------
          |
      ---------
          |
      ---------
          |
      ---------
          |
      ---------
          |
      ---------
          |
      ---------
          |
      ---------
          |
      ---------
 3.3V [2] |
      ---------
 MOSI [4] |
      ---------
 MISO [7] |
      ---------
          | [3] CE0
      ---------
  GND [1] |
      ---------

      ------------
      3.5 JACK | O
      ------------
```
Mapping
- 3 Debug Clock -> 27 / 10
- 4 Debug Data -> 28 / 12
- 7 Reset 24 / 13

0. Download firmware

!!! source_routing - max 5 devices !!!

```
cd flash_cc2531
wget https://github.com/Koenkk/Z-Stack-firmware/raw/master/coordinator/Z-Stack_Home_1.2/bin/source_routing/CC2531_SOURCE_ROUTING_20190619.zip
rm -rf CC2531ZNP*
unzip CC2531_SOURCE_ROUTING_20190619.zip
```
1. Check
```
./cc_chipid -c 10 -d 12 -r 13
  ID = b524.
```
2. Erase
```
./cc_erase  -c 10 -d 12 -r 13
  ID = b524.
  erase result = 00a2.
```
3. Flash
```
./cc_write -c 10 -d 12 -r 13 CC2531ZNP-Prod.hex
  ID = b524.
  reading line 15490.
  file loaded (15497 lines read).
writing page 128/128.
verifying page 128/128.
 flash OK.
```

#### Wake on LAN
```
sudo ethtool -s INTERFACE
```
```
sudo vim /etc/systemd/system/wol.service
```
```
[Unit]
Description=Configure Wake On LAN

[Service]
Type=oneshot
ExecStart=/sbin/ethtool -s INTERFACE wol g

[Install]
WantedBy=basic.target
```
```
sudo systemctl daemon-reload
sudo systemctl enable wol.service
sudo systemctl start wol.service
```

#### Wemos D1 mini
- Connect WEMOS to USB-cable
- Place a lead between D3 and GND
- Press the small reset button
- Remove lead to D3
- You *should* be able upload firmware