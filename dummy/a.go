package main

import (
	"log"
	"net"
)

func pnc(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	const myAddr = ":4567"
	conn, err := net.ListenPacket("udp", myAddr)
	pnc(err)
	defer conn.Close()

	log.Printf("(udp server) started on %+v", myAddr)
	ackRequest := make([]byte, 100)
	n, addr, err := conn.ReadFrom(ackRequest[:])
	pnc(err)
	log.Printf("(udp server) ack req from  router[%v]", string(ackRequest[:n]))
	_, err = conn.WriteTo([]byte("ack\n"), addr)
	pnc(err)
	log.Printf("(udp server) acknowledged  router[%v]", string(ackRequest[:n]))
}
