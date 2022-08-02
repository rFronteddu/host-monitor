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
	m := &measure.Measure{
		Subject:  measure.Subject_os,
		Strings:  make(map[string]string),
		Integers: make(map[string]int64),
		Doubles:  make(map[string]float64),
	}
	fmt.Printf("Sending report to %s, %s\n", udpc.destination, "Hello")
	udpc.send([]byte("Hello"))
	ticker := time.NewTimer(udpc.period)
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
