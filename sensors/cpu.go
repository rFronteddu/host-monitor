package sensors

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"hostmonitor/measure"
	"time"
)

type CPU struct {
	period time.Duration
}

func NewCPUSensor(period time.Duration) *CPU {
	sensor := new(CPU)
	sensor.period = period
	return sensor
}

func (sensor *CPU) Poll(measure *measure.Measure) {
	// needs no sleep since getting the cpu will take time
	v, _ := cpu.Percent(sensor.period, false)
	fmt.Printf("CPU Report %s AVG Used CPU Percent: %f%%\n", sensor.period, v[0])
	measure.Integers["CPU_AVG"] = int64(v[0])
}
