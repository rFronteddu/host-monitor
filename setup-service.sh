#! /bin/bash
# This script will set up (and re-set-up) the daemon for the host-monitor

# build the main program
go build main.go

# copy all the necessary program files for the daemon
sudo rm -R /usr/lib/monitor
sudo mkdir /usr/lib/monitor
sudo cp host-monitor.service /usr/lib/monitor/host-monitor.service
sudo cp hostmonitor.yaml /usr/lib/monitor/hostmonitor.yaml
sudo cp main /usr/lib/monitor/main

# copy the service to the
sudo rm /etc/systemd/system/host-monitor.service
sudo cp host-monitor.service /etc/systemd/system/host-monitor.service
sudo systemctl enable host-monitor

# kill previously running programs and start
sudo pkill main
sudo systemctl start host-monitor

# view daemon status at anytime with >sudo systemctl status host-monitor