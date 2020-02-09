package example

import (
	pb "ledstripinterface/proto"
	"log"
	"sync"
	"time"
)

func Play(render func(frame pb.Frame) error) {
	var lastState state
	stateUpdated := false
	m := sync.Mutex{}

	go func() {
		state := NewState(216)
		transmitTicker := time.NewTicker(time.Second / 100)
		for {
			select {
			case <-transmitTicker.C:
				state.Update()
				m.Lock()
				lastState = *state
				stateUpdated = true
				m.Unlock()
			}
		}
	}()

	const frameDuration = time.Second / 100 // Hz
	for {
		startTime := time.Now()
		var frame pb.Frame
		renderNecessary := false
		m.Lock()
		if stateUpdated {
			frame = lastState.ToFrame()
			stateUpdated = false
			renderNecessary = true
		}
		m.Unlock()
		if renderNecessary {
			err := render(frame)
			if err != nil {
				log.Print(err)
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

type state struct {
	baseColor  pb.Color
	brightness float32
	upwards    bool
	minimum    float32
	maximum    float32
	timePassed time.Duration
	pixels     []float32
	backBuffer []float32
}

func NewState(nLed int) *state {
	makeCenter := func() []float32 {
		a := make([]float32, nLed)
		for i := range a {
			a[i] = 0.5
		}
		return a
	}
	return &state{
		// baseColor: pb.Color{
		// 	R: 255,
		// 	G: 180,
		// 	B: 255,
		// },
		// warm white, similar to halogen, determined by experiment with support of https://encycolorpedia.com
		baseColor: pb.Color{
			R: 233,
			G: 130,
			B: 35,
		},
		brightness: 0.5,
		upwards:    true,
		minimum:    0.5,
		maximum:    0.99,
		timePassed: 0,
		pixels:     makeCenter(),
		backBuffer: makeCenter(),
	}
}

func (s *state) Update() {
	// global breathe
	if s.brightness >= s.maximum {
		s.upwards = false
	}
	if s.brightness <= s.minimum {
		s.upwards = true
	}
	if s.upwards {
		s.brightness += 0.0075
	} else {
		s.brightness -= 0.0075
	}

	for i := range s.pixels {
		s.pixels[i] = s.brightness
	}

	/*
		// randomly distribute
		s.pixels[rand.Int()%len(s.pixels)] = s.brightness
	*/

	/*
		// sparkle
		for i := range s.pixels {
			s.pixels[i] = 0
		}
		for i := 0; i < len(s.pixels)/10; i++ {
			s.pixels[rand.Int()%len(s.pixels)] = 1
		}
	*/

	/*
		// random glow
		s.timePassed += time.Millisecond * 10
		if s.timePassed > time.Second {
			s.timePassed = 0
			const size = 10
			for i := size / 2; i < len(s.pixels)-size/2; i += size {
				dst := s.minimum + (rand.Float32() / (s.maximum + s.minimum))
				for j := i - size/2; j < i+size/2; j++ {
					s.backBuffer[j] = dst
				}
			}
		}
		for i := range s.pixels {
			s.pixels[i] += (s.backBuffer[i] - s.pixels[i]) / 50

		}
	*/
}

func (s *state) ToFrame() pb.Frame {
	frame := pb.Frame{}
	for _, pixel := range s.pixels {
		color := pb.Color{
			R: uint32(pixel * float32(s.baseColor.R)),
			G: uint32(pixel * float32(s.baseColor.G)),
			B: uint32(pixel * float32(s.baseColor.B)),
		}
		frame.Pixels = append(frame.Pixels, &color)
	}
	return frame
}
