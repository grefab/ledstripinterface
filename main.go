package main

import (
	"ledstripinterface/arduino"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)

	// init and wait for response
	comPort := "COM3"
	err := arduino.EstablishConnection(comPort)
	if err != nil {
		log.Fatal(err)
	}

	// display something fancy with stable fps
	const frameDuration time.Duration = time.Second / 10
	const num = 16
	var i uint8 = 0
	for {
		startTime := time.Now()
		arduino.SendStrip(arduino.MakeStrip(num, i))
		i = (i + 1) % num
		duration := time.Now().Sub(startTime)

		pause := frameDuration - duration
		if pause <= 0 {
			log.Printf("render overflow")
		} else {
			time.Sleep(pause)
		}
	}
}
