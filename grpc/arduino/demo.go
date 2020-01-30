package arduino

import (
	"github.com/golang/protobuf/proto"
	pb "ledstripinterface/pb"
	"log"
	"time"
)

func Breathe(render func(frame pb.Frame) error) {
	const frameDuration = time.Second / 500 // Hz
	brightness := 0.5

	upwards := true
	minimum := 0.3
	maximum := 1.0

	lastFrame := makeFullFrame(pb.Color{R: 128, G: 128, B: 128})
	for {
		startTime := time.Now()
		{
			frame := makeFullFrame(pb.Color{
				R: uint32(233 * brightness),
				G: uint32(130 * brightness),
				B: uint32(35 * brightness),
			})
			if !proto.Equal(&frame, &lastFrame) {
				lastFrame = frame
				err := render(frame)
				if err != nil {
					log.Print(err)
				}
			}
			if brightness >= maximum {
				upwards = false
			}
			if brightness <= minimum {
				upwards = true
			}
			if upwards {
				brightness += 0.01
			} else {
				brightness -= 0.01
			}
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

func makeFullFrame(color pb.Color) pb.Frame {
	nLed := 216
	frame := pb.Frame{}
	for i := 0; i < nLed; i++ {
		frame.Pixels = append(frame.Pixels, &color)
	}
	return frame
}

func RunDemo(send func(frame []*pb.Color) error) {
	frames := make(chan []*pb.Color, 100)
	go updateState(frames)

	// display something fancy with stable fps
	const frameDuration = time.Second / 10
	for {
		pixels := <-frames
		startTime := time.Now()
		err := send(pixels)
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
	for i := 0; i < 23; i++ {
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

			var intensity uint32 = 16
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

		time.Sleep(time.Millisecond * 1500)
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
