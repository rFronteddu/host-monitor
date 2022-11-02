package sensors

import (
	"github.com/shirou/gopsutil/v3/net"
	"hostmonitor/measure"
	"log"
	"time"
)

type NetSensor struct {
	period time.Duration
}

func NewNetSensor(period time.Duration) *NetSensor {
	nets := new(NetSensor)
	nets.period = period
	return nets
}

func (nets *NetSensor) Poll(measure *measure.Measure) {
	time.Sleep(nets.period)
	iList, _ := net.Interfaces()
	// if pernic is true, returns result divided by interface, returns a summary otherwise
	counter, _ := net.IOCounters(false)
	log.Printf("Net Report - Interfaces: %v Counters: %v\n", iList, counter)
}
