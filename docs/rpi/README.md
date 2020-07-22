## Отключение swap

```
sudo dphys-swapfile swapoff
sudo dphys-swapfile uninstall
sudo update-rc.d dphys-swapfile remove
sudo apt purge dphys-swapfile
```

## Обновление

```
sudo apt update
sudo apt full-upgrade
sudo apt -y dist-upgrade
sudo rpi-update
sudo reboot
```
