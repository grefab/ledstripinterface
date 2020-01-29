package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"ledstripinterface/grpc/arduino"
	pb "ledstripinterface/pb"
	"net"
)

type displayServer struct {
	controller arduino.Controller
}

func RunServer(ip net.IP, port int, comPort string) {
	d := displayServer{}
	err := d.controller.EstablishConnection(comPort)
	if err != nil {
		panic(err)
	}

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
	d.controller.SendStrip(frame.Pixels)
	return &empty.Empty{}, nil
}
