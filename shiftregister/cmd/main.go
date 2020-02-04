package main

import (
	pb "ledstripinterface/proto"
	"ledstripinterface/service"
	"ledstripinterface/shiftregister"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)

	// We have a conveyor prism carrying a vial every 38.1mm, a total length of our display LED strip of 1500mm and
	// 216 LEDs over that length. That means we have 5.4864 LEDs per prism. For 39 prisms that results in 213.97 LEDs
	// which we round up to 214.
	sr := shiftregister.ShiftRegister{LedCount: 214}
	for i := 0; i < 38; i++ {
		sr.Add(pb.Color{
			R: 255,
			G: 255,
			B: 255,
		})
	}
	for len(sr.Vials) < 39 {
		sr.Add(pb.Color{R: 0, G: 0, B: 0})
	}
	log.Print("running in client mode")
	display := service.NewRemoteDisplay("localhost:15050")

	framePipe := make(chan pb.Frame, 100)
	go func() {
		const frameDuration = time.Second / 100 // Hz
		for {
			frame := <-framePipe
			startTime := time.Now()
			err := display.ShowFrame(&frame)
			if err != nil {
				log.Print(err)
			}
			duration := time.Now().Sub(startTime)
			pause := frameDuration - duration
			if pause <= 0 {
				log.Printf("render overflow: %v", pause)
			} else {
				time.Sleep(pause)
			}
		}
	}()

	for {
		white := pb.Color{
			R: 255,
			G: 255,
			B: 255,
		}

		startTime := time.Now()
		frames := 0
		sr.Shift(white, func(frame pb.Frame) { framePipe <- frame; frames++ })
		log.Printf("frames: %v, duration: %v", frames, time.Now().Sub(startTime))
		time.Sleep(time.Millisecond * 500)
	}
}
