```
sudo cp -f /home/pi/go/src/github.com/kihamo/boggart/cmd/boggart/boggart.service /lib/systemd/system/boggart.service
sudo systemctl daemon-reload
sudo journalctl -f -u boggart.service
```

## Server
#### First
```
gox -output="cmd/server/boggart" -osarch="linux/amd64" ./cmd/server/
sudo cp -f /home/kihamo/cmd/server/boggart.service /lib/systemd/system/boggart.service
sudo cp -f /home/kihamo/cmd/server/boggart /usr/local/bin/boggart-server
sudo chmod +x /usr/local/bin/boggart-server
sudo systemctl daemon-reload
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service
```
#### Update
```
gox -output="cmd/server/boggart" -osarch="linux/amd64" ./cmd/server/
sudo cp -f /home/kihamo/cmd/server/boggart /usr/local/bin/boggart-server
sudo chmod +x /usr/local/bin/boggart-server
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
```

## Agent
#### First
```
GOARM=7 gox -output="cmd/agent/boggart" -osarch="linux/arm" ./cmd/agent/
sudo cp -f /home/pi/cmd/agent/boggart.service /lib/systemd/system/boggart.service
sudo cp -f /home/pi/cmd/agent/boggart /usr/local/bin/boggart-agent
sudo chmod +x /usr/local/bin/boggart-agent
sudo systemctl daemon-reload
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service
```
#### Update
```
GOARM=7 gox -output="cmd/agent/boggart" -osarch="linux/arm" ./cmd/agent/
sudo cp -f /home/pi/cmd/agent/boggart /usr/local/bin/boggart-agent
sudo chmod +x /usr/local/bin/boggart-agent
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
```

## Other

```
sudo systemctl stop boggart.service
sudo systemctl restart boggart.service && sudo journalctl -f -u boggart.service
sudo systemctl start boggart.service && sudo journalctl -f -u boggart.service

sudo echo "22" > /sys/class/gpio/unexport
sudo echo "27" > /sys/class/gpio/unexport
```