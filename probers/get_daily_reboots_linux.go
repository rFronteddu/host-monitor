//go:build linux || darwin

package probers

import (
	"bytes"
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
	var cmd = exec.Command("bash", "./reboots.sh")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error trying to get # of reboots: %v\n", err)
		return 100
	}
	time.Sleep(10 * time.Millisecond)
	reboots, err := strconv.Atoi(strings.Split(out.String(), "\n")[0])
	if err != nil {
		fmt.Printf("Error parsing system response: %v\n", err)
		return 100
	}
	//fmt.Printf("Number of reboots today: %v\n", reboots)
	return int64(reboots)
}
