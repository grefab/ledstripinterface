package main

import (
	"flag"
	"ledstripinterface/demo"
	pb "ledstripinterface/proto"
	"ledstripinterface/service"
	"ledstripinterface/udp"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	clientMode := flag.Bool("client", false, "start in client mode, tests connection to server, needs endpoint")
	endpoint := flag.String("endpoint", "0.0.0.0:15050", "endpoint to listen to in server mode or to connect to in client mode")
	udpTarget := flag.String("udpTarget", "10.42.0.57:1337", "IP and port for UDP target")
	flag.Parse()

	if *clientMode {
		log.Print("running in client mode")
		display := service.NewRemoteDisplay(*endpoint)
		demo.PlayShiftRegisterDemo(func(frame pb.Frame) error {
			return display.ShowFrame(&frame)
		})
	} else {
		log.Print("running in server mode. see -help for other options")
		ip, port := service.ParseEndpoint(*endpoint)
		log.Printf("running gRPC service on %v:%v to supply UPD on arduino: %v", ip, port, *udpTarget)
		service.RunServer(ip, port, func(frame pb.Frame) {
			data := make([]byte, 0, len(frame.Pixels)*3)
			for _, color := range frame.Pixels {
				data = append(data, byte(color.R))
				data = append(data, byte(color.G))
				data = append(data, byte(color.B))
			}
			err := udp.Send(*udpTarget, data)
			if err != nil {
				log.Println(err)
			}
		})
	}
}
