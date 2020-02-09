package main

import (
	pb "ledstripinterface/proto"
	"ledstripinterface/service"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Print("running in client mode")
	display := service.NewRemoteDisplay("localhost:15050")

	srRed := pb.ShiftRegister{
		Vials:  nil,
		Offset: 0,
		Stride: 1,
	}
	for i := 0; i < 20; i++ {
		srRed.Vials = append(srRed.Vials, &pb.Color{
			R: 255,
			G: 0,
			B: 0,
		})
	}
	srGreen := pb.ShiftRegister{
		Vials:  nil,
		Offset: 1,
		Stride: 1,
	}
	for i := 0; i < 20; i++ {
		srGreen.Vials = append(srGreen.Vials, &pb.Color{
			R: 0,
			G: 255,
			B: 0,
		})
	}

	conveyor := pb.Conveyor{
		Strip: &pb.Strip{
			LengthMeters:       1.5,
			LedCount:           216,
			ChainElementSizeMm: 38.1,
		},
	}
	conveyor.ShiftRegisters = append(conveyor.ShiftRegisters, &srRed)
	conveyor.ShiftRegisters = append(conveyor.ShiftRegisters, &srGreen)

	err := display.ShowConveyor(&conveyor)
	if err != nil {
		log.Print(err)
	}

	for {
		startTime := time.Now()

		err = display.Move()
		if err != nil {
			log.Print(err)
		}

		time.Sleep(time.Millisecond * 1000)

		err := display.ShowConveyor(&conveyor)
		if err != nil {
			log.Print(err)
		}

		time.Sleep(time.Millisecond * 1000)
		log.Printf("duration: %v", time.Now().Sub(startTime))
	}
}
