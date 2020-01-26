package main

import (
	"flag"
	"ledstripinterface/lib"
	"log"
)

// func main() {
// 	go func() {
// 		ip, port := lib.ParseEndpoint("127.0.0.1:15050")
// 		lib.RunServer(ip, port, "COM3")
// 	}()
//
// 	display := lib.NewRemoteDisplay("127.0.0.1:15050")
//
// 	// display something fancy with stable fps
// 	const frameDuration = time.Second / 10
// 	const num = 16
// 	var i uint32 = 0
// 	for {
// 		startTime := time.Now()
// 		err := display.ShowFrame(&pb.Frame{
// 			Frames: arduino.MakeStrip(num, i),
// 		})
// 		if err != nil {
// 			panic(err)
// 		}
// 		i = (i + 1) % num
// 		duration := time.Now().Sub(startTime)
//
// 		pause := frameDuration - duration
// 		if pause <= 0 {
// 			log.Printf("render overflow: %v", pause)
// 		} else {
// 			time.Sleep(pause)
// 		}
// 	}
// }

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds)
	endpoint := flag.String("endpoint", "0.0.0.0:15050", "endpoint to listen to in server mode or to connect to in client mode")
	serialPort := flag.String("serialPort", "/dev/ttyUSB0", "serial device to connect to, e.g. 'COM3' for windows")
	flag.Parse()
	ip, port := lib.ParseEndpoint(*endpoint)
	log.Printf("rs232 to arduino: %v, running gRPC server on %v:%v", *serialPort, ip, port)
	lib.RunServer(ip, port, *serialPort)
}
