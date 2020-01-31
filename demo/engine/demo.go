package engine

import (
	pb "ledstripinterface/pb"
	"log"
	"math/rand"
	"sync"
	"time"
)

func Play(render func(frame pb.Frame) error) {
	var lastState state
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
				m.Unlock()
			}
		}
	}()

	const frameDuration = time.Second / 100 // Hz
	for {
		startTime := time.Now()
		m.Lock()
		frame := lastState.ToFrame()
		m.Unlock()
		err := render(frame)
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

type state struct {
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
		brightness: 0.5,
		upwards:    true,
		minimum:    0.3,
		maximum:    1.0,
		timePassed: 0,
		pixels:     makeCenter(),
		backBuffer: makeCenter(),
	}
}

func (s *state) Update() {
	/*
		// global breathe
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
	*/

	/*
		// randomly distribute
		s.pixels[rand.Int()%len(s.pixels)] = s.brightness
	*/
	/*
		// sparkle
		for i := range s.pixels {
			s.pixels[i] = 0
		}
		for i := 0; i < len(s.pixels) / 10; i++ {
			s.pixels[rand.Int()%len(s.pixels)] = 1
		}
	*/

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
