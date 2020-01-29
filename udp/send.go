package udp

import (
	"log"
	"net"
)

func Send(addr string, data []byte) {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	log.Print(a.IP, a.Port)
	conn, err := net.DialUDP("udp", nil, a)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	n, err := conn.Write(data)
	log.Print(n)
}
