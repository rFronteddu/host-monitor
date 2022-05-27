package sensors

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"hostmonitor/measure"
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
	fmt.Printf("Disk Report - Usage: %v\n", usage)
	measure.Integers["DISK_USAGE"] = int64(usage.UsedPercent)
	measure.Integers["DISK_FREE"] = int64(usage.Free)
}
