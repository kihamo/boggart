#! /bin/sh

for id in core sdram_c sdram_i sdram_p ; do
	echo "cpu_voltage,id=$id value=$(vcgencmd measure_volts $id | tr -d 'volt=V')"
done