package engine

import (
	pb "ledstripinterface/pb"
	"log"
	"time"
)

func Play(render func(frame pb.Frame) error) {
	stateUpdates := make(chan state, 10)
	go func() {
		state := NewState(216)
		for {
			select {
			case <-time.After(time.Millisecond * 10):
				state.Update()
				stateUpdates <- *state
			}
		}
	}()

	var lastState state
	for {
		select {
		case state := <-stateUpdates:
			lastState = state
		case <-time.After(time.Millisecond * 10):
			frame := lastState.ToFrame()
			err := render(frame)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

type state struct {
	brightness float32
	upwards    bool
	minimum    float32
	maximum    float32
	pixels     []float32
}

func NewState(nLed int) *state {
	return &state{
		brightness: 0.5,
		upwards:    true,
		minimum:    0.3,
		maximum:    1.0,
		pixels:     make([]float32, nLed),
	}
}

func (s *state) Update() {
	if s.brightness >= s.maximum {
		s.upwards = false
	}
	if s.brightness <= s.minimum {
		s.upwards = true
	}
	if s.upwards {
		s.brightness += 0.01
	} else {
		s.brightness -= 0.01
	}

	for i := range s.pixels {
		s.pixels[i] = s.brightness
	}

	// randomly distribute
	// s.pixels[rand.Int()%len(s.pixels)] = s.brightness
}

func (s *state) ToFrame() pb.Frame {
	frame := pb.Frame{}
	for _, pixel := range s.pixels {
		color := pb.Color{
			R: uint32(233 * pixel),
			G: uint32(130 * pixel),
			B: uint32(35 * pixel),
		}
		frame.Pixels = append(frame.Pixels, &color)
	}
	return frame
}
