package sensors

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	"hostmonitor/measure"
	"time"
)

type Host struct {
	period time.Duration
}

var lastBootTime = uint64(time.Now().Unix())
var dailyRebootCounter int64 = 0

func NewHostSensor(period time.Duration) *Host {
	sensor := new(Host)
	sensor.period = period
	return sensor
}

func (sensor *Host) Poll(measure *measure.Measure) {
	time.Sleep(sensor.period)
	// Reset the reboot counter every midnight
	if time.Now().Hour() == 0 && time.Now().Minute() == 0 {
		dailyRebootCounter = 0
	}
	// Get host info
	h, _ := host.Info()
	if lastBootTime < h.BootTime {
		lastBootTime = h.BootTime
		dailyRebootCounter++
	}
	// Print report and send to Measure channel
	fmt.Printf("Host Report - Host ID: %v Host Name: %v, OS: %v, Platform: %v, Arch: %v, Boot Date: %v, Reboots Today: %v, Uptime: %v\n", h.HostID, h.Hostname, h.OS, h.Platform, h.KernelArch, time.Unix(int64(h.BootTime), 0), dailyRebootCounter, h.Uptime)
	measure.Strings["host_id"] = h.HostID
	measure.Strings["host_name"] = h.Hostname
	measure.Strings["os"] = h.OS
	measure.Strings["platform"] = h.Platform
	measure.Strings["kernelArch"] = h.KernelArch
	measure.Strings["bootTime"] = time.Unix(int64(h.BootTime), 0).Format(time.RFC822)
	measure.Integers["Reboots_Today"] = dailyRebootCounter
}
