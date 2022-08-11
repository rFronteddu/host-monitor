package transport

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/proto"
	"hostmonitor/measure"
	"net"
	"time"
)

type UDPClient struct {
	measureCh   chan *measure.Measure
	destination string
	period      time.Duration
}

func NewUDPClient(destination string, inCh chan *measure.Measure, period time.Duration) *UDPClient {
	udpc := new(UDPClient)
	udpc.measureCh = inCh
	udpc.destination = destination
	udpc.period = period
	return udpc
}

func (udpc *UDPClient) Start() {
	go udpc.sendReports()
}

func (udpc *UDPClient) sendReports() {
	// Send out initial handshake
	m0 := &measure.Measure{
		Subject:  measure.Subject_os,
		Strings:  make(map[string]string),
		Integers: make(map[string]int64),
		Doubles:  make(map[string]float64),
	}
	m0.Strings["host_id"] = "Hello"
	out0, marshalError0 := proto.Marshal(m0)
	if marshalError0 != nil {
		fmt.Printf("Encountered an error while marshaling measure %v\n", marshalError0)
	}
	fmt.Printf("Sending handshake to %s, %s\n", udpc.destination, m0.String())
	udpc.send(out0)

	ticker := time.NewTicker(udpc.period)
	m := &measure.Measure{
		Subject:  measure.Subject_os,
		Strings:  make(map[string]string),
		Integers: make(map[string]int64),
		Doubles:  make(map[string]float64),
	}
	for {
		select {
		case <-ticker.C:
			m.Timestamp = &timestamp.Timestamp{Seconds: time.Now().Unix()}
			out, marshalError := proto.Marshal(m)
			if marshalError != nil {
				fmt.Printf("Encountered an error while marshaling measure %v\n", marshalError)
				break
			}
			fmt.Printf("Sending report to %s, %s\n", udpc.destination, m.String())
			udpc.send(out)
			m = &measure.Measure{
				Subject:  measure.Subject_os,
				Strings:  make(map[string]string),
				Integers: make(map[string]int64),
				Doubles:  make(map[string]float64),
			}
		case msg := <-udpc.measureCh:
			// copy all values
			for k, v := range msg.Doubles {
				m.Doubles[k] = v
			}
			for k, v := range msg.Integers {
				m.Integers[k] = v
			}
			for k, v := range msg.Strings {
				m.Strings[k] = v
			}
			if m.Integers["uptime"] != 0 && m.Integers["uptime"] < int64(udpc.period.Seconds()) {
				m.Integers["reboots"] = 1
			}
		}
	}
}

func (udpc *UDPClient) send(buff []byte) {
	conn, err := net.Dial("udp", udpc.destination)
	if err != nil {
		fmt.Printf("Sending buffer returned an error: %v\n", err)
		return
	}
	_, writeError := conn.Write(buff)
	if writeError != nil {
		fmt.Printf("conn.Write(m) failed with error: %v\n", writeError)
	}
}
