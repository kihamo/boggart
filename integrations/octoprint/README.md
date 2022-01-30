## Install OctoPi

## Upgrade

```
sudo rpi-update
sudo apt-get full-upgrade
```

## Translate
https://github.com/AlexanderPro29/OctoPrint-RU-LangPack/releases/tag/1.7.2

## Install Webcam

#### edit /boot/octopi.txt
```
camera="usb"
additional_brokenfps_usb_devices=("046d:0825")
```

#### UDEV static link for Logitech, Inc. Webcam C270

```
cat > /etc/udev/rules.d/01-webcam-usb.rules
SUBSYSTEMS=="usb", ATTR{idVendor}=="046d", ATTR{idProduct}=="0825", ATTRS{serial}=="112FCD10", SYMLINK+="LogitechC270"

sudo udevadm control --reload-rules && udevadm trigger

ls -la /dev/LogitechC270
```