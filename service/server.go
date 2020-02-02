package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	pb "ledstripinterface/proto"
	"net"
)

type displayServer struct {
	SendFrame func(frame pb.Frame)
}

func RunServer(ip net.IP, port int, sendFrame func(frame pb.Frame)) {
	d := displayServer{SendFrame: sendFrame}
	s := grpc.NewServer()
	pb.RegisterDisplayServer(s, &d)

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{IP: ip, Port: port})
	if err != nil {
		panic(err)
	}
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func (d *displayServer) ShowFrame(_ context.Context, frame *pb.Frame) (*empty.Empty, error) {
	d.SendFrame(*frame)
	return &empty.Empty{}, nil
}
