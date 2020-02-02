package main

import (
	"ledstripinterface/demo/engine"
	pb "ledstripinterface/proto"
	"ledstripinterface/serial"
	"ledstripinterface/udp"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	render := func(frame pb.Frame) error {
		return udp.Send("10.42.0.57:1337", serial.FrameToBytes(frame))
	}
	engine.Play(render)
}
