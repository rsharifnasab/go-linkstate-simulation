package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
)

type Router struct {
	conn       *net.UDPConn
	connReader *bufio.Reader
	connWriter *bufio.Writer

	port    int
	index   int
	logFile *os.File

	routersCount     int
	neighbours       []*Edge
	netConns         [][]*Edge
	forwardingTable  map[int]int
	shortestPathTree []*Edge

	portMap           map[int]int
	mergedPortMaps    map[int]int
	mergedPortMapLock sync.RWMutex

	managerConnection *net.TCPConn
	managerWriter     *bufio.Writer
	managerReader     *bufio.Reader

	forwardingChannel  chan string
	managerPacketsDone chan struct{}
}

type Edge struct {
	Dest int
	Cost int
}

type Packet struct {
	source      int
	destination int
	data        string
}

func NewRouter() *Router {
	return &Router{
		managerPacketsDone: make(chan struct{}),
		mergedPortMapLock:  sync.RWMutex{},
	}
}

func (router *Router) sendPacket(rawPacket string) {
	packet := parsePacket(rawPacket)
	if packet.destination == router.index {
		log.Printf("Receive packet from #%v. [%v]", packet.source, packet.data)
	} else {
		nextHop, ok := router.forwardingTable[packet.destination]
		if !ok {
			log.Printf("problem forwarding packet (%v) from router #%v to %v\n",
				rawPacket, packet.source, packet.destination)
			return
		}
		log.Printf("forwarding [%v] to nextHop router #%v\n", rawPacket, nextHop)
		router.writeUDPAsBytes(nextHop, []byte(rawPacket))
	}
}

func (router *Router) forwardPacketsFromManager() {
	for {
		data := router.readStringFromManager()
		if data == "QUIT" {
			close(router.managerPacketsDone)
		} else {
			//log.Printf("recieved  [%v] from manager\n", data)
			router.sendPacket(data)
		}
	}
}

// determine packets in udp buffer
// they are from previous broadcasting
func isPacketBad(packet string) bool {
	return packet[0] == '{'
}

func (router *Router) forwardPacketsFromOtherRouters() {
	for {
		packet := string(router.readUDPAsBytes())
		for isPacketBad(packet) {
			packet = string(router.readUDPAsBytes())
		}
		router.sendPacket(packet)
	}
}
