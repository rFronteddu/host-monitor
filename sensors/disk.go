package sensors

import (
	"github.com/shirou/gopsutil/v3/disk"
	"hostmonitor/measure"
	"log"
	"time"
)

type Disk struct {
	period time.Duration
}

func NewDiskSensor(period time.Duration) *Disk {
	sensor := new(Disk)
	sensor.period = period
	return sensor
}

func (sensor *Disk) Poll(measure *measure.Measure) {
	time.Sleep(sensor.period)
	usage, _ := disk.Usage("/")
	log.Printf("Disk Report - Usage: %v\n", usage)
	measure.Integers["diskUsedPercent"] = int64(usage.UsedPercent)
	measure.Integers["diskFree"] = int64(usage.Free) / 1024 / 1024 / 1024
	measure.Integers["diskUsed"] = int64(usage.Used) / 1024 / 1024 / 1024
	measure.Integers["diskTotal"] = int64(usage.Total) / 1024 / 1024 / 1024
}
