package arduino

import (
	"fmt"
	"hostmonitor/task"
	"time"
)

type Monitor struct {
}

func NewArduinoMonitor() *Monitor {
	arduinoMonitor := new(Monitor)
	return arduinoMonitor
}

func (monitor *Monitor) Start(BOARD_IP string) {
	fmt.Printf("Starting Arduino Monitor...\n")
	ticker := time.NewTicker(60 * time.Minute)
	lastArduinoReachableTimestamp := time.Time{}
	go func() {
		for _ = range ticker.C {
			println("\nPinging Arduino...\n")
			pingTask := task.NewPingTask()
			if pingTask.Start(BOARD_IP) == true {
				lastArduinoReachableTimestamp = time.Now()
			}
			fmt.Printf("\nLast Arduino reachable timestamp: %s", lastArduinoReachableTimestamp.String())
		}
	}()
}
