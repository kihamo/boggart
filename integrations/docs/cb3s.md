1. Docs
   2. https://docs.libretiny.eu/docs/flashing/tools/ltchiptool/
   3. https://docs.libretiny.eu/docs/platform/beken-72xx/
   3. https://pypi.org/project/ltchiptool/3.0.0a3/
   4. https://docs.libretiny.eu/boards/cb3s/
   5. https://esphome.io/components/libretiny.html
2. Install cli `sudo pip install ltchiptool`
2. Dump factory firmware `ltchiptool flash read -d /dev/cu.usbserial-1420 bk7231n dump.bin`
   Уou need to bridge **CEN** pin to **GND** with a wire.
3. Flash firmware `ltchiptool flash write -d /dev/cu.usbserial-1420 firmware.uf2`