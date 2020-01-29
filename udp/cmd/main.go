package main

import (
	"ledstripinterface/grpc/arduino"
	pb "ledstripinterface/pb"
	"ledstripinterface/udp"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	arduino.RunDemo(func(frame []*pb.Color) error {
		return udp.Send("10.42.0.57:1337", arduino.StripToBytes(frame))
	})
}
