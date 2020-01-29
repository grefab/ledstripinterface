package main

import (
	"flag"
	"ledstripinterface/grpc"
	"ledstripinterface/grpc/arduino"
	pb "ledstripinterface/pb"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	clientMode := flag.Bool("client", false, "start in client mode, tests connection to server, needs endpoint")
	endpoint := flag.String("endpoint", "0.0.0.0:15050", "endpoint to listen to in server mode or to connect to in client mode")
	serialPort := flag.String("serialPort", "/dev/ttyUSB0", "serial device to connect to, e.g. 'COM3' for windows")
	flag.Parse()

	if *clientMode {
		log.Print("running in client mode")
		display := grpc.NewRemoteDisplay(*endpoint)
		arduino.RunDemo(func(frame []*pb.Color) error {
			return display.ShowFrame(&pb.Frame{
				Pixels: frame,
			})
		})
	} else {
		log.Print("running in server mode. see -help for other options")
		ip, port := grpc.ParseEndpoint(*endpoint)
		log.Printf("rs232 to arduino: %v, running gRPC server on %v:%v", *serialPort, ip, port)
		grpc.RunServer(ip, port, *serialPort)
	}
}
