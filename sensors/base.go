package sensors

import (
	"fmt"
	"hostmonitor/measure"
	"time"
)

type BaseSensor struct {
	sensor    sensor
	measureCh chan<- *measure.Measure
	tag       string
}

func NewSensor(sensor sensor, tag string, outCh chan *measure.Measure) *BaseSensor {
	base := new(BaseSensor)
	base.sensor = sensor
	base.tag = tag
	base.measureCh = outCh
	return base
}

func (base *BaseSensor) Start() {
	fmt.Printf("\tStarting %s\n", base.tag)
	go func() {
		for {
			base.poll()
		}
	}()
}

func (base *BaseSensor) poll() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal error: %s, %s will sleep 5 seconds before attempting to proceed\n", err, base.tag)
		}
		time.Sleep(5 * time.Second)
	}()
	m := &measure.Measure{
		Strings:  make(map[string]string),
		Integers: make(map[string]int64),
		Doubles:  make(map[string]float64),
	}
	base.sensor.Poll(m)
	base.measureCh <- m
}
