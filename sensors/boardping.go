package sensors

import (
	"fmt"
	"hostmonitor/measure"
	"hostmonitor/probers"
	"time"
)

type boardPing struct {
	period       time.Duration
	boardMonitor probers.Monitor
}

func NewBoardPingSensor(period time.Duration, boardMonitor *probers.Monitor) *boardPing {
	sensor := new(boardPing)
	sensor.period = period
	sensor.boardMonitor = *boardMonitor
	return sensor
}

func (sensor *boardPing) Poll(measure *measure.Measure) {
	fmt.Printf("Last reached board at %s\n", probers.GetLastReachableTimestamp(&sensor.boardMonitor))
	measure.Integers["LastArduinoReachableTimestamp"] = int64(time.Now().Sub(probers.GetLastReachableTimestamp(&sensor.boardMonitor)).Seconds())
}
