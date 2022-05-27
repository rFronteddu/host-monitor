package grpc

import (
	"context"
	"google.golang.org/grpc"
	pb "hostmonitor/pinger"
	"hostmonitor/probers"
	"log"
	"net"
	"strconv"
)

type PingerProxy struct {
	port string
}

func NewPingerProxy(port int) *PingerProxy {
	receiver := new(PingerProxy)
	receiver.port = ":" + strconv.Itoa(port)
	return receiver
}

func (proxy *PingerProxy) Start() {
	go func() {
		for {
			proxy.manageCommands()
		}
	}()
}

// server is used to implement pinger.PingServer.
type server struct {
	pb.UnimplementedPingerServer
}

// Ping implements pinger.PingServer
func (s *server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	log.Printf("Received ping request towards: %v", in.GetTargetAddress())

	replyChannel := make(chan *pb.PingReply)

	icmp := probers.NewICMPProbe(in.GetTargetAddress(), replyChannel)
	icmp.Start()
	reply := <-replyChannel
	return reply, nil
}

func (proxy *PingerProxy) manageCommands() {
	lis, err := net.Listen("tcp", proxy.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPingerServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
