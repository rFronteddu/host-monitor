package sensors

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/load"
	"hostmonitor/measure"
	"time"
)

type Load struct {
	period time.Duration
}

func NewLoadSensor(period time.Duration) *Load {
	sensor := new(Load)
	sensor.period = period
	return sensor
}

func (sensor *Load) Poll(measure *measure.Measure) {
	time.Sleep(sensor.period)
	loadAvg, _ := load.Avg()
	fmt.Printf("Load Report - %v\n", loadAvg)
	measure.Integers["LOAD_1"] = int64(loadAvg.Load1)
	measure.Integers["LOAD_5"] = int64(loadAvg.Load5)
	measure.Integers["LOAD_15"] = int64(loadAvg.Load15)
}
