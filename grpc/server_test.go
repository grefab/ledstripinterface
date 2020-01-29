package grpc

import (
	"ledstripinterface/arduino"
	pb "ledstripinterface/pb"
	"log"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	go func() {
		ip, port := ParseEndpoint("127.0.0.1:15050")
		RunServer(ip, port, "COM3")
	}()

	display := NewRemoteDisplay("127.0.0.1:15050")

	// display something fancy with stable fps
	const frameDuration = time.Second / 10
	const num = 16
	var i uint32 = 0
	for {
		startTime := time.Now()
		err := display.ShowFrame(&pb.Frame{
			Frames: arduino.MakeStrip(num, i),
		})
		if err != nil {
			panic(err)
		}
		i = (i + 1) % num
		duration := time.Now().Sub(startTime)

		pause := frameDuration - duration
		if pause <= 0 {
			log.Printf("render overflow: %v", pause)
		} else {
			time.Sleep(pause)
		}
	}
}
