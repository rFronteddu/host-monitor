package probers

import (
	"fmt"
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

func (monitor *Monitor) Start(BOARD_IP string) {
	fmt.Printf("Starting periodic pinger...\n")
	replyCh := make(chan *pb.PingReply)
	var p *pb.PingReply
	ticker := time.NewTicker(60 * time.Second)
	monitor.lastReachable = time.Time{}
	go func() {
		for _ = range ticker.C {
			fmt.Printf("\nPinging board @ %s...\n", BOARD_IP)
			icmp := NewICMPProbe(BOARD_IP, replyCh)
			icmp.Start()
			p = <-replyCh
			if p.Reachable == true {
				monitor.lastReachable = time.Now()
				fmt.Printf("\nBoard %s was reached at %s", BOARD_IP, monitor.lastReachable.String())
			} else {
				fmt.Printf("\nBoard %s is unreachable!", BOARD_IP)
			}
		}
	}()
}

func GetLastReachableTimestamp(boardMonitor *Monitor) time.Time {
	return boardMonitor.lastReachable
}
