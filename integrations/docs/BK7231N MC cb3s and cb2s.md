CB3S
=====
1. Docs
   2. https://docs.libretiny.eu/docs/flashing/tools/ltchiptool/
   3. https://docs.libretiny.eu/docs/platform/beken-72xx/
   3. https://pypi.org/project/ltchiptool/3.0.0a3/
   4. https://docs.libretiny.eu/boards/cb3s/
   5. https://esphome.io/components/libretiny.html
2. Install cli `sudo pip install ltchiptool`
3. Dump factory firmware `ltchiptool flash read -d /dev/cu.usbserial-1420 bk7231n dump.bin`
   Уou need to bridge **CEN** pin to **GND** with a wire.
4. Flash firmware `ltchiptool flash write -d /dev/cu.usbserial-1420 firmware.uf2`

CB2S
=====

ltchiptool -v flash write -d /dev/cu.usbserial-1420 ~/Downloads/villa-garage-plug.uf2

Docs: https://docs.libretiny.eu/boards/cb2s/?h=cb2s#pinout

Перед включением `ltchiptool flash write` скинуть RX и TX, дождаться инструкции как подключить и накинуть контакты 
RX и TX только после этого хоть какая то загрузка пошла. Эффект повторить не удалось на том же устройстве, но прошивка залилась. 
На новом устройстве удалось воспроизвести.

Как и у CB3S у чипа такая же проблема -- слабый wifi, подключается к точке но через 10-15 скидывается и так пока не стабилизируется.

See CB3S but:
```
   RX -> RX1
   TX -> TX1
```