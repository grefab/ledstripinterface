package main

import (
	"flag"
	"ledstripinterface/demo"
	pb "ledstripinterface/proto"
	"ledstripinterface/serial"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	serialPort := flag.String("serialPort", "/dev/ttyUSB0", "serial device to connect to, e.g. 'COM3' for windows")
	flag.Parse()

	log.Print("sending demo data to serial port", *serialPort)
	controller := serial.Controller{}
	err := controller.EstablishConnection(*serialPort)
	if err != nil {
		panic(err)
	}
	demo.PlayShiftRegisterDemo(func(frame pb.Frame) error {
		controller.SendFrame(frame)
		return nil
	})
}
