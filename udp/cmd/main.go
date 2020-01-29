package main

import "ledstripinterface/udp"

func main() {
	buf := []byte("ping")
	for {
		udp.Send("10.42.0.57:1337", buf)
	}
}
