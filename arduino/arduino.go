package arduino

import (
	"github.com/tarm/serial"
	"log"
	"strings"
	"time"
)

type Color struct {
	R uint8
	G uint8
	B uint8
}

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

	done := make(chan bool)
	eval := func(cmd string) {
		switch cmd {
		case "READY":
			log.Println("arduino signals ready")
			done <- true
		case "FLUSH":
			// log.Println("flushed")
		default:
			log.Printf("unknown arduino message: %v", cmd)
		}
	}
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
				eval(splits[0])
				splits = splits[1:]
			}
			received = splits[0]
		}
	}
	go listenToSerial()
	<-done
	return nil
}

func SendStrip(strip []Color) {
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
		data := []byte{color.R, color.G, color.B}
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
