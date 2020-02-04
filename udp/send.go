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

	// dim data so we do not stress LEDs too much. we measured 2.2A for 216 LEDs (all white).
	// no dimming caused the first third to be bright, the rest of the strip getting dimmer then.
	// no significant temperature increase with this setting
	// note: dimming to 25% is done by ardiuno program
	// for i := range data {
	// 	data[i] /= 4
	// }
	_, err = conn.Write(data)
	return nil
}
