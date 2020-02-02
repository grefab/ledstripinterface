package main

import (
	"flag"
	"ledstripinterface/demo"
	"ledstripinterface/grpc"
	pb "ledstripinterface/proto"
	"ledstripinterface/serial"
	"ledstripinterface/udp"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	clientMode := flag.Bool("client", false, "start in client mode, tests connection to server, needs endpoint")
	endpoint := flag.String("endpoint", "0.0.0.0:15050", "endpoint to listen to in server mode or to connect to in client mode")
	serialPort := flag.String("updTarget", "10.13.37.10:1337", "IP and port for UDP target")
	flag.Parse()

	if *clientMode {
		log.Print("running in client mode")
		display := grpc.NewRemoteDisplay(*endpoint)
		demo.PlayShiftRegisterDemo(func(frame pb.Frame) error {
			return display.ShowFrame(&frame)
		})
	} else {
		log.Print("running in server mode. see -help for other options")
		ip, port := grpc.ParseEndpoint(*endpoint)
		log.Printf("gRPC to UPD on arduino: %v, running gRPC server on %v:%v", *serialPort, ip, port)
		grpc.RunServer(ip, port, func(frame pb.Frame) {
			err := udp.Send(*serialPort, serial.FrameToBytes(frame))
			if err != nil {
				log.Println(err)
			}
		})
	}
}
