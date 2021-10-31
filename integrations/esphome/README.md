# ESPHome Dashboard install and run

```
cd /Users/kihamo/go/src/github.com/kihamo/boggart/integrations/esphome
docker pull esphome/esphome
docker run --rm -p 6052:6052 -p 6123:6123 -v "${PWD}":/config -e ESPHOME_DASHBOARD_USE_PING=true -it esphome/esphome
```