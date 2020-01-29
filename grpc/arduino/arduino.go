package arduino

import (
	"github.com/tarm/serial"
	pb "ledstripinterface/pb"
	"log"
	"strings"
)

type Controller struct {
	conn     *serial.Port
	received chan string
}

func (controller *Controller) EstablishConnection(comPort string) error {
	conn, err := serial.OpenPort(&serial.Config{
		Name: comPort,
		Baud: 115200,
	})
	if err != nil {
		return err
	}
	controller.conn = conn
	controller.received = make(chan string, 10)
	listenToSerial := func() {
		var received string
		for {
			buf := make([]byte, 128)
			n, err := conn.Read(buf)

			if err != nil {
				log.Fatal(err)
			}
			received += string(buf[:n])
			// log.Printf("received: %v\n", received)
			splits := strings.Split(received, "\n")
			// log.Printf("splits: %v\n", splits)

			// When we split and no separator is present, we get an
			// array of size 1 containing the source string.
			// When we split and an a separator is present at the end,
			// we get an array with 2 (or more) elements, with
			// the last element is empty. If the last element is not
			// empty, it means that command is not completely received yet.
			for len(splits) > 1 {
				controller.received <- splits[0]
				splits = splits[1:]
			}
			received = splits[0]
		}
	}
	go listenToSerial()

	return nil
}

func (controller *Controller) SendStrip(strip []*pb.Color) {
	data := StripToBytes(strip)
	data = append(data, 0)
	data = append(data, 0)
	data = append(data, 0)
	controller.write(data)
	<-controller.received // should be OK
}

func StripToBytes(strip []*pb.Color) []byte {
	data := make([]byte, 0, len(strip)*3)
	for _, color := range strip {
		// Our protocol treats 000 as flush, so we exchange color values of 0 by 1 here.
		// Should have no visible effect for LEDs.
		if color.R == 0 {
			color.R = 1
		}
		if color.G == 0 {
			color.G = 1
		}
		if color.B == 0 {
			color.B = 1
		}
		data = append(data, byte(color.R))
		data = append(data, byte(color.G))
		data = append(data, byte(color.B))
	}
	return data
}

func (controller *Controller) write(data []byte) {
	_, err := controller.conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func MakeStrip(num int, highlight uint32) (strip []*pb.Color) {
	for i := 0; i < num; i++ {
		strip = append(strip, &pb.Color{
			R: highlight * 3,
			G: highlight * 3,
			B: highlight * 3,
		})
	}

	strip[highlight] = &pb.Color{
		R: 255,
		G: 255,
		B: 255,
	}
	return strip
}
