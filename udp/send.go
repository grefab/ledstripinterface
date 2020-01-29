package udp

import (
	"net"
)

func Send(addr string, data []byte) error {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, a)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	return nil
}
