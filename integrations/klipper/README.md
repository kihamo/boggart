## Upgrade

```
cd ~/klipper/
git pull
make clean
make

sudo service klipper stop
make flash FLASH_DEVICE=/dev/serial/by-id/usb-FTDI_FT232R_USB_UART_AD0JMICN-if00-port0
sudo service klipper start
```

## Undervoltage warnings
Решение тут https://community.octoprint.org/t/put-tape-on-the-5v-pin-why-and-how/13574
RPi начинает питать собой плату управления принтером и поэтому проседает по вольтажу. Хардкорный костыль -- изолировать контакт 5V в проводе от RPi до платы принтера