package lib

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"ledstripinterface/arduino"
	pb "ledstripinterface/pb"
	"net"
)

type server struct {
}

func RunServer(ip net.IP, port int, comPort string) {
	err := arduino.EstablishConnection(comPort)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterDisplayServer(s, &server{})

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{IP: ip, Port: port})
	if err != nil {
		panic(err)
	}
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func (s *server) ShowFrame(_ context.Context, frame *pb.Frame) (*empty.Empty, error) {
	arduino.SendStrip(frame.Pixels)
	return &empty.Empty{}, nil
}
