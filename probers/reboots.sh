#!/bin/bash
#This script used only by linux/darwin
day=$(date | awk {'print $1 " " $3 " " $2'})
reboot_times=$(last reboot | grep "$day" | wc -l)
echo $reboot_times
