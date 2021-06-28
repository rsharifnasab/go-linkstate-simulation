package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func (router *Router) testNeighbouringLinks() {
	log.Printf("(ack client) checking neighbouring links started")
	for _, edge := range router.neighbours {
		index := edge.Dest
		port := router.portMap[index]
		//log.Printf("dialing to router[%v] on port %v\n", index, port)
		conn := dialUDP(fmt.Sprintf("localhost:%v", port))
		router.sendAckRequest(conn, index)
		router.getAckResponse(conn, index)
		conn.Close()
		//log.Printf("%v check", edge.Dest)
	}
	router.writeToManager("ACKS_RECEIVED")
	log.Printf("(ack client) checking neighbouring links done")
}

func (router *Router) sendAckRequest(conn net.Conn, index int) {
	writer := bufio.NewWriter(conn)
	writer.WriteString(fmt.Sprintf("%v\n", router.index))
	writer.Flush()
	log.Printf("(ack client) send ack request to router #%v", index)
}

func (router *Router) getAckResponse(conn net.Conn, index int) {
	ackResponse, err := bufio.NewReader(conn).ReadString('\n')
	pnc(err)
	if ackResponse != "ack\n" {
		log.Fatal(
			fmt.Sprintf(
				"(ack client) router #%v  by saying %v didn't ack",
				index, ackResponse))
	}
	log.Printf("(ack client) received ack from router #%v", index)

}

func (router *Router) sendAcknowledgements() {
	log.Printf("(ack server) listening for other routers check started")
	for i := 0; i < len(router.neighbours); i++ {
		ackRequest := make([]byte, 100)
		n, addr, err := router.conn.ReadFrom(ackRequest[:])
		pnc(err)
		//log.Printf("(udp server ack) ack req from  router[%v]", string(ackRequest[:n-1]))
		router.conn.WriteTo([]byte("ack\n"), addr)
		log.Printf("(ack server) acknowledged router #%v", string(ackRequest[:n-1]))
	}
}

func (router *Router) waitNetworkReadiness() {
	str := router.readStringFromManager()
	if str != "NETWORK_READY" {
		log.Fatal("Manager didn't sent network ready")
	} else {
		log.Printf("Network is ready")
	}
}
