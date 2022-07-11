#!/bin/bash
#windows
day=$(Get-Date)
log=$(Get-WinEvent -FilterHashtable @{LogName = 'System';id=6006; StartTime=$day})
"$log".Count

#linux
#day=$(date | awk {'print $1 " " $3 " " $2'})
#reboot_times=$(last reboot | select -first 10 | grep "$day" | wc -l)
#echo $reboot_times
