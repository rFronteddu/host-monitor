package sensors

import (
	"github.com/shirou/gopsutil/v3/load"
	"hostmonitor/measure"
	"log"
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
	log.Printf("Load Report - %v\n", loadAvg)
	measure.Integers["load1"] = int64(loadAvg.Load1)
	measure.Integers["load5"] = int64(loadAvg.Load5)
	measure.Integers["load15"] = int64(loadAvg.Load15)
}
