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
	ticker := time.NewTicker(10 * time.Second)
	//lastReached := time.Time{}
	//dailyUnreachableCount := 0
	go func() {
		for _ = range ticker.C {
			println("Pinging Arduino") ///////////////////////
			pingTask := task.NewPingTask()
			pingTask.Start(BOARD_IP)
		}
	}()
}
