package sensors

import (
	"github.com/shirou/gopsutil/v3/mem"
	"hostmonitor/measure"
	"log"
	"time"
)

type VirtualMemorySensor struct {
	period time.Duration
}

func NewVirtualMemorySensor(period time.Duration) *VirtualMemorySensor {
	vms := new(VirtualMemorySensor)
	vms.period = period
	return vms
}

func (vms *VirtualMemorySensor) Poll(measure *measure.Measure) {
	time.Sleep(vms.period)
	v, _ := mem.VirtualMemory()
	log.Printf("Memory Report - Total: %v, Free: %v, UsedPercent: %f%%\n", v.Total, v.Free, v.UsedPercent)
	measure.Integers["vmUsedPercent"] = int64(v.UsedPercent)
	measure.Integers["vmFree"] = int64(v.Free) / 1024 / 1024
	measure.Integers["vmUsed"] = int64(v.Used) / 1024 / 1024
	measure.Integers["vmTotal"] = int64(v.Total) / 1024 / 1024
}
