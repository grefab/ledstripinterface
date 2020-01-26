package arduino

import (
	"github.com/tarm/serial"
	pb "ledstripinterface/pb"
	"log"
	"time"
)

var connection *serial.Port

func EstablishConnection(comPort string) error {
	conn, err := serial.OpenPort(&serial.Config{
		Name: comPort,
		Baud: 57600,
	})
	if err != nil {
		return err
	}
	connection = conn
	return nil
}

func SendStrip(strip []*pb.Color) {
	startTime := time.Now()
	for _, color := range strip {
		// Our protocol treats 000 as flush,
		// so we exchange 0s by 1s here.
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
		data := []byte{byte(color.R), byte(color.G), byte(color.B)}
		write(data)
	}
	write([]byte{0, 0, 0})
	duration := time.Now().Sub(startTime)

	// Communication takes time. At 57600 baud (fastest that tested stable) we send 16 LEDs every ~8ms.
	// Upon flush, Arduino needs to transfer data over the WS2812 bus, which also needs to be taken into account,
	// otherwise its serial buffer overflows. This time is a ratio of RS232 transmission time and has been
	// determined experimentally.
	time.Sleep(duration / 6)
}

func write(data []byte) {
	written := 0
	for written < len(data) {
		n, err := connection.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		written += n
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
