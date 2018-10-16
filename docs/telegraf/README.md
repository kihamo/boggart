## Установка и настройка telegraf

```
https://dl.influxdata.com/telegraf/releases/telegraf_1.8.1-1_armhf.deb
dpkg -i telegraf_1.8.1-1_armhf.deb && rm -rf telegraf_1.8.1-1_armhf.deb

cat > /etc/telegraf/telegraf.conf
```

```
cat > /etc/telegraf/cpu_frequentie.sh

#! /bin/sh

NUM_PROC=$(cat /sys/devices/system/cpu/present | sed 's/-/ /g' | awk '{print $2}')

for i in $(seq 0 $NUM_PROC); do
	echo "cpu_frequentie,cpu=cpu$i value=$(cat /sys/devices/system/cpu/cpu$i/cpufreq/scaling_cur_freq)"
done

chmod +x /etc/telegraf/cpu_frequentie.sh
```

```
cat > /etc/telegraf/cpu_voltage.sh

#! /bin/sh

for id in core sdram_c sdram_i sdram_p ; do
	echo "cpu_voltage,id=$id value=$(vcgencmd measure_volts $id | tr -d 'volt=V')"
done

chmod +x /etc/telegraf/cpu_voltage.sh
```

## Чтобы заработал cpu_voltage.sh

```
sudo usermod -G video telegraf
```