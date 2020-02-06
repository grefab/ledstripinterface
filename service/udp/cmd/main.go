package main

import (
	"ledstripinterface/example/engine"
	pb "ledstripinterface/proto"
	"ledstripinterface/service/serial"
	udp2 "ledstripinterface/service/udp"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	render := func(frame pb.Frame) error {
		return udp2.Send("10.42.0.57:1337", serial.FrameToBytes(frame))
	}
	engine.Play(render)
}
