#! /bin/sh

NUM_PROC=$(cat /sys/devices/system/cpu/present | sed 's/-/ /g' | awk '{print $2}')

for i in $(seq 0 $NUM_PROC); do
	echo "cpu_frequentie,cpu=cpu$i value=$(cat /sys/devices/system/cpu/cpu$i/cpufreq/scaling_cur_freq)"
done