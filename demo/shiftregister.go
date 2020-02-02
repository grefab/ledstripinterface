package demo

import (
	pb "ledstripinterface/proto"
	"log"
	"time"
)

func PlayShiftRegisterDemo(send func(frame pb.Frame) error) {
	frames := make(chan []*pb.Color, 100)
	go updateState(frames)

	// display something fancy with stable fps
	const frameDuration = time.Second / 10
	for {
		pixels := <-frames
		startTime := time.Now()
		err := send(pb.Frame{Pixels: pixels})
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
}

type Vial struct {
	Color pb.Color
}
type Conveyor struct {
	Vials []Vial
}

func (c *Conveyor) fill() {
	pattern := "gggb...gggbgggg"
	for i := 0; i < 60; i++ {
		v := pattern[i%len(pattern)]
		switch v {
		case '.':
			c.Vials = append(c.Vials,
				Vial{
					Color: pb.Color{
						R: 16,
						G: 16,
						B: 16,
					}})
		case 'b':
			c.Vials = append(c.Vials,
				Vial{
					Color: pb.Color{
						R: 255,
						G: 0,
						B: 0,
					}})
		default:
			c.Vials = append(c.Vials,
				Vial{
					Color: pb.Color{
						R: 255,
						G: 255,
						B: 255,
					}})
		}
	}
}
func updateState(colors chan []*pb.Color) {
	future := Conveyor{}
	future.fill()
	conveyor := Conveyor{}
	conveyor.fill()

	for {
		var strip []*pb.Color
		addColor := func(color pb.Color) { strip = append(strip, &color) }
		addBlack := func() { addColor(pb.Color{R: 0, G: 0, B: 0}) }
		sendStrip := func() {
			var bufStrip []*pb.Color
			for _, e := range strip {
				color := *e
				bufStrip = append(bufStrip, &color)
			}

			var intensity uint32 = 255
			bufStrip[4] = &pb.Color{R: intensity, G: intensity, B: 0}
			bufStrip[7] = &pb.Color{R: intensity, G: intensity, B: 0}
			colors <- bufStrip
		}

		for _, vial := range conveyor.Vials {
			addBlack()
			addColor(vial.Color)
			addColor(vial.Color)
			addBlack()
		}
		sendStrip()

		time.Sleep(time.Millisecond * 500)
		if len(future.Vials) == 0 {
			future.fill()
		}

		// transition to right
		strip = strip[1:]
		addBlack()
		sendStrip()

		strip = strip[1:]
		addColor(future.Vials[0].Color)
		sendStrip()

		strip = strip[1:]
		addColor(future.Vials[0].Color)
		sendStrip()

		conveyor.Vials = conveyor.Vials[1:]
		conveyor.Vials = append(conveyor.Vials, Vial{Color: future.Vials[0].Color})
		future.Vials = future.Vials[1:]
	}
}
