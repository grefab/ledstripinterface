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
	sr := pb.ShiftRegister{LedCount: 214}
	for i := 0; i < 39; i++ {
		shiftregister.Add(&sr, pb.Color{
			R: 255,
			G: 255,
			B: 255,
		})
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
		red := pb.Color{
			R: 255,
			G: 0,
			B: 0,
		}

		startTime := time.Now()
		frames := shiftregister.Shift(red, &sr)
		log.Printf("frames: %v, duration: %v", len(frames), time.Now().Sub(startTime))

		for _, frame := range frames {
			framePipe <- frame
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
