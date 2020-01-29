package main

import (
	"ledstripinterface/grpc/arduino"
	pb "ledstripinterface/pb"
	"ledstripinterface/udp"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	render := func(frame []*pb.Color) error {
		return udp.Send("10.42.0.57:1337", arduino.StripToBytes(frame))
	}
	arduino.Breathe(render)
}
