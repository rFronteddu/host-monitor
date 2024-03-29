package sensors

import (
	"github.com/shirou/gopsutil/v3/host"
	"hostmonitor/measure"
	"log"
	"time"
)

type Host struct {
	period time.Duration
}

func NewHostSensor(period time.Duration) *Host {
	sensor := new(Host)
	sensor.period = period
	return sensor
}

func (sensor *Host) Poll(measure *measure.Measure) {
	time.Sleep(sensor.period)

	h, _ := host.Info()
	log.Printf("Host Report - Host ID: %v Host Name: %v, OS: %v, Platform: %v, Arch: %v, Boot Date: %v, Uptime: %v\n", h.HostID, h.Hostname, h.OS, h.Platform, h.KernelArch, time.Unix(int64(h.BootTime), 0), h.Uptime)

	measure.Strings["hostId"] = h.HostID
	measure.Strings["hostname"] = h.Hostname
	measure.Strings["os"] = h.OS
	measure.Strings["platform"] = h.Platform
	measure.Strings["kernelArch"] = h.KernelArch
	measure.Integers["bootTime"] = int64(h.BootTime)
	measure.Integers["uptime"] = int64(h.Uptime)
}
