# ESPHome Dashboard install and run

```
cd /Users/kihamo/go/src/github.com/kihamo/boggart/integrations/esphome
docker pull esphome/esphome
docker run --rm -p 6052:6052 -p 6123:6123 -v "${PWD}":/config -e ESPHOME_DASHBOARD_USE_PING=true -it esphome/esphome
```

# Wemos D1 mini ESP8266
## Flash
- Connect WEMOS to USB-cable
- Place a lead between D3 and GND
- Press the small reset button
- Remove lead to D3
- You *should* be able upload firmware

# NodeMCU ESP32
## Flash
- Power ON
- Press and keep BOOT button
- You *should* be able upload firmware

# Live D1 mini ESP32
## Pinout
![Pinout](docs/MH-ET_LIVE_D1_mini_ESP32_pinout.png)

## Flashing
Magic :)
- Power ON
- Press Flash ESP in esphome-flasher
- Power Reset button