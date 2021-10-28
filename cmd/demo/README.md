## Files structure
`demo` - compiled application file

`boggart.env` - setting environment variables to start the service

`boggart.service` - configuration file for the service on the systemd

`config.yaml` - binds (devices) configuration file

## Install
#### 1. Build release
```
gox -output="cmd/server/demo" -osarch="linux/amd64" -ldflags="-X 'main.Name=Boggart Demo' -X 'main.Version=`date +"%y%m%d"`' -X 'main.Build=`date +"%H%M%S"`'" ./cmd/demo/
```

#### 2. Upload release to host system
Upload binary files `demo`, `boggart.env`, `boggart.service` and `config.yaml` to directory `/opt/boggart/` on host system and copy executable file `demo`
```
sudo cp -f /opt/boggart/boggart /usr/local/bin/boggart-demo
sudo chmod +x /usr/local/bin/boggart-demo
```

#### 3. Initial steps on host system
Copy configs
```
sudo cp -f /opt/boggart/boggart.service /lib/systemd/system/boggart.service
```

Install boggart as service
```
sudo cp -f /opt/boggart/boggart.service /lib/systemd/system/boggart.service
sudo systemctl daemon-reload
sudo systemctl enable boggart.service
```

#### 3. Start service
```
sudo systemctl start boggart.service
```

#### 4. Check install
You can show log
```
sudo journalctl -f -u boggart.service
```