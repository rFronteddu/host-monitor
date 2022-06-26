package probers

import (
	"fmt"
	"hostmonitor/measure"
	pb "hostmonitor/pinger"
	"time"
)

type Monitor struct {
	lastReachable time.Time
}

func NewBoardMonitor() *Monitor {
	boardMonitor := new(Monitor)
	boardMonitor.lastReachable = time.Time{}
	return boardMonitor
}

func (monitor *Monitor) Start(BOARD_IP string, inCh chan *measure.Measure) {
	fmt.Printf("Starting periodic pinger...\n")
	replyCh := make(chan *pb.PingReply)
	var p *pb.PingReply
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for _ = range ticker.C {
			fmt.Printf("\nPinging board @ %s...\n", BOARD_IP)
			icmp := NewICMPProbe(BOARD_IP, replyCh)
			icmp.Start()
			p = <-replyCh
			if p.Reachable == true {
				monitor.lastReachable = time.Now()
				fmt.Printf("Board %s was reached at %s\n", BOARD_IP, monitor.lastReachable.String())
				m := &measure.Measure{
					Strings:  make(map[string]string),
					Integers: make(map[string]int64),
					Doubles:  make(map[string]float64),
				}
				m.Integers["LastArduinoReachableTimestamp"] = int64(time.Now().Sub(monitor.lastReachable).Seconds())
				inCh <- m
			} else {
				fmt.Printf("\nBoard %s is unreachable!", BOARD_IP)
			}
		}
	}()
}
