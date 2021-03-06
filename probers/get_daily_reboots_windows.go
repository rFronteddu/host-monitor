//go:build windows

package probers

import (
	"fmt"
	"hostmonitor/measure"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type RebootCounter struct {
	reboots int64
}

func NewRebootCounter() *RebootCounter {
	rebootCounter := new(RebootCounter)
	rebootCounter.reboots = 0
	return rebootCounter
}

func (rebootCounter *RebootCounter) Start(inCh chan *measure.Measure) {
	fmt.Printf("Starting reboot counter...\n")
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for _ = range ticker.C {
			rebootCounter.reboots = GetReboots()
			m := &measure.Measure{
				Strings:  make(map[string]string),
				Integers: make(map[string]int64),
				Doubles:  make(map[string]float64),
			}
			m.Integers["Reboots_Today"] = rebootCounter.reboots
			inCh <- m
		}
	}()
}

// GetReboots This function runs reboots.sh which returns the number of reboots for the current day
func GetReboots() int64 {
	c1 := "day=$(Get-Date)"
	c2 := "log=$(Get-WinEvent -FilterHashtable @{LogName = 'System';id=6006; StartTime=$day})"
	c3 := "\"$log\".Count"
	cmd := exec.Command("powershell.exe", "-c", "chcp", "65001", ">", "$null", ";", c1)
	cmd = exec.Command("powershell.exe", "-c", "chcp", "65001", ">", "$null", ";", c2)
	cmd = exec.Command("powershell.exe", "-c", "chcp", "65001", ">", "$null", ";", c3)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error trying to get # of reboots: %v\n", err)
		return 100
	}
	reboots, err := strconv.Atoi(strings.Split(string(out), "\r\n")[0])
	if err != nil {
		fmt.Printf("Error parsing system response: %v\n", err)
		return 100
	}
	//fmt.Printf("Number of reboots today: %v\n", reboots)
	return int64(reboots)
}
