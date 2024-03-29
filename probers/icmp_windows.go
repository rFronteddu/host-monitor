//go:build windows

package probers

import (
	"bytes"
	pb "hostmonitor/pinger"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func ping(target string, replyCh chan *pb.PingReply) {
	var cmd = exec.Command("ping", target)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Printf("Unable to ping %s\n", target)
		replyCh <- &pb.PingReply{
			Reachable:      false,
			AvgRtt:         0,
			LostPercentage: 0,
		}
		return
	}

	var packetsSent int
	var packetsReceived int
	var rtt int

	s1 := strings.Split(out.String(), "\n")
	for _, s := range s1 {
		log.Println(s)
		if strings.Contains(s, "unreachable") {
			log.Printf("Destination %s unreachable\n", target)
			replyCh <- &pb.PingReply{
				Reachable:      false,
				AvgRtt:         0,
				LostPercentage: 100,
			}
			return
		}
		if strings.Contains(s, "Packets: Sent") {
			packetsSent, packetsReceived = extractReceivedAndSent(s)
		}
		if strings.Contains(s, "Minimum = ") {
			rtt = extractAverageRTT(s)
		}
	}

	log.Printf("Packet sent: %v Packets received: %v Avg rtt: %v\n", packetsSent, packetsReceived, rtt)
	replyCh <- &pb.PingReply{
		Reachable:      true,
		AvgRtt:         int32(rtt),
		LostPercentage: int32(100 - (100 * packetsSent / packetsReceived)),
	}

}

func extractAverageRTT(out string) int {
	var avgRTT int
	s2 := strings.Split(out, ",")
	s3 := strings.Replace(s2[2], "Average = ", "", -1)
	s4 := strings.Replace(s3, "ms", "", -1)
	avgRTT, _ = strconv.Atoi(strings.TrimSpace(s4))
	return avgRTT
}

func extractReceivedAndSent(out string) (int, int) {
	var packetsSent int
	var packetsReceived int

	s2 := strings.Split(out, ",")
	s3 := strings.Replace(s2[0], "Packets: Sent =", "", -1)
	s4 := strings.TrimSpace(s3)
	packetsSent, _ = strconv.Atoi(s4)

	s5 := strings.Replace(s2[1], "Received =", "", -1)
	s6 := strings.TrimSpace(s5)
	packetsReceived, _ = strconv.Atoi(s6)

	return packetsSent, packetsReceived
}
