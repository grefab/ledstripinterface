package lib

import (
	pb "ledstripinterface/proto"
	"log"
	"net"
	"regexp"
	"strconv"
)

func ParseEndpoint(endpoint string) (net.IP, int) {
	parse := regexp.MustCompile(`^(.*):(.*)$`)
	matches := parse.FindAllStringSubmatch(endpoint, -1)
	ip := net.ParseIP(matches[0][1])
	if ip == nil {
		log.Fatalf("not a valid ip: %v", matches[0][1])
	}
	port, err := strconv.Atoi(matches[0][2])
	if err != nil {
		log.Fatalf("not a valid port: %v: %v", matches[0][2], err)
	}
	return ip, port
}

func FrameToBytes(frame *pb.Frame) []byte {
	data := make([]byte, 0, len(frame.Pixels)*3)
	for _, color := range frame.Pixels {
		data = append(data, byte(color.R))
		data = append(data, byte(color.G))
		data = append(data, byte(color.B))
	}
	return data
}
