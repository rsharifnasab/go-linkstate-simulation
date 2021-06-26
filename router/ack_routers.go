package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func (router *Router) testNeighbouringLinks() {
	log.Printf("(ack client) checking neighbouring links {{")
	for _, edge := range router.neighbours {
		index := edge.Dest
		port := router.portMap[index]
		//log.Printf("dialing to router[%v] on port %v\n", index, port)
		conn := dialUDP(fmt.Sprintf("localhost:%v", port))
		router.sendAckRequest(conn, index, port)
		router.getAckResponse(conn, index, port)
		conn.Close()
		//log.Printf("%v check", edge.Dest)
	}
	router.writeToManager("ACKS_RECEIVED")
	log.Printf("(ack client) }}")
}

func (router *Router) sendAckRequest(conn net.Conn, index, port int) {
	writer := bufio.NewWriter(conn)
	writer.WriteString(fmt.Sprintf("%v\n", router.index))
	writer.Flush()
	log.Printf("(ack client) sent ack request to %v on %v\n", index, port)
}

func (router *Router) getAckResponse(conn net.Conn, index, port int) {
	ackResponse, err := bufio.NewReader(conn).ReadString('\n')
	pnc(err)
	if ackResponse != "ack\n" {
		log.Fatal(
			fmt.Sprintf(
				"(ack client) router #%v on port %v by saying %v didn't ack",
				index, port, ackResponse))
	}
	log.Printf("(ack client) received ack from %v on %v\n", index, port)

}

func (router *Router) sendAcknowledgements() {
	log.Printf("(ack server) listening for other routers")
	for i := 0; i < len(router.neighbours); i++ {
		ackRequest := make([]byte, 100)
		n, addr, err := router.conn.ReadFrom(ackRequest[:])
		pnc(err)
		//log.Printf("(udp server ack) ack req from  router[%v]", string(ackRequest[:n-1]))
		router.conn.WriteTo([]byte("ack\n"), addr)
		log.Printf("(ack server) acknowledged  router[%v]", string(ackRequest[:n-1]))
	}
}

func (router *Router) waitNetworkReadiness() {
	str := router.readStringFromManager()
	if str != "NETWORK_READY" {
		log.Fatal("manager didn't sent network ready")
	} else {
		log.Printf("network is ready")
	}
}
