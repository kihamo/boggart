# NodeMCU

## I2C

I2C can be used to connect up to 127 nodes via a bus that only requires two data wires, known as SDA and SCL.

- SDA=4 => D2.
- SCL=5 => D1

```
default I2C Ports of <Wire.h>

GPIO 5  - D1 is SDA
GPIO 4  - D2 is SCL

NodeMCU - D1 ( GPIO 5  - D1 - SCL)
NodeMCU - D2 ( GPIO 4  - D2 - SDA)
NodeMCU - D5 ( GPIO 14 - D5 )
NodeMCU - D6 ( GPIO 12 - D6 )
NodeMCU - D7 ( GPIO 13 - D7 )
NodeMCU - D8 ( GPIO 15 - D8 )
```

## SPI
SPI is much simpler than I2C. Master and slave are linked by three data wires, usually called MISO, (Master in, Slave out), MOSI (Master out, Slave in) and M-CLK.

- M-CLK => D5
- MISO => D6
- MOSI => D7

(SPI Bus SS (CS)is D8.)

## Links

- https://github.com/andyprv/workshop/blob/master/Demo/bme280test/bme280test.ino

## NodeMCU

esptool.py --port=/dev/cu.SLAB_USBtoUART write_flash -fm=dio -fs=4MB 0x00000 ~/web-server/firmware/adc,bit,cron,crypto,file,gpio,http,i2c,mqtt,net,node,ow,rotary,rtctime,sjson,spi,tmr,uart,wifi-integer.bin

