package main

import (
	"flag"
	pb "ledstripinterface/proto"
	"ledstripinterface/service/lib"
	"ledstripinterface/udp"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	endpoint := flag.String("endpoint", "0.0.0.0:15050", "endpoint to listen to in server mode or to connect to in client mode")
	udpTarget := flag.String("udpTarget", "10.42.0.57:1337", "IP and port for UDP target")
	flag.Parse()

	ip, port := lib.ParseEndpoint(*endpoint)
	log.Printf("running gRPC service on %v:%v to supply UPD on arduino: %v", ip, port, *udpTarget)
	lib.RunServer(ip, port, func(frame pb.Frame) {
		data := lib.FrameToBytes(&frame)
		err := udp.Send(*udpTarget, data)
		if err != nil {
			log.Println(err)
		}
	})
}
