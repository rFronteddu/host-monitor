package sensors

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"hostmonitor/measure"
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
	fmt.Printf("Memory Report - Total: %v, Free: %v, UsedPercent: %f%%\n", v.Total, v.Free, v.UsedPercent)
	measure.Integers["vm_used_percent"] = int64(v.UsedPercent)
	measure.Integers["vm_free"] = int64(v.Free / 1024 / 1024)
}
