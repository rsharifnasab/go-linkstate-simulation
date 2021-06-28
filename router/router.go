package main

import (
	"bufio"
	"fmt"
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
	From      int
	To        int
	CostToNow int

	Data string
}

func NewRouter() *Router {
	return &Router{
		managerPacketsDone: make(chan struct{}),
		mergedPortMapLock:  sync.RWMutex{},
	}
}

func (packet *Packet) serialize() []byte {
	return []byte(fmt.Sprintf("%d %d %d %s",
		packet.From, packet.To, packet.CostToNow, packet.Data))
}
func (router *Router) writePacketTo(nextHop int, packet *Packet) {
	router.writeUDPAsBytes(nextHop, packet.serialize())
}

func (router *Router) addEdgeCostTo(nextHop int, packet *Packet) {
	for _, v := range router.neighbours {
		if v.Dest == nextHop {
			packet.CostToNow += v.Cost
			return
		}
	}
	log.Fatalf("cannot find next hop cost (%v) for packet (%+v)", nextHop, packet)
}
func (router *Router) sendPacket(rawPacket string) {

	packet := parsePacket(rawPacket)
	if packet.To == router.index {
		log.Printf("Receive packet from #%v with cost %v. [%v]", packet.From, packet.CostToNow, packet.Data)
	} else {
		nextHop, ok := router.forwardingTable[packet.To]
		if !ok {
			log.Printf("problem forwarding packet (%v) from router #%v to %v\n",
				rawPacket, packet.From, packet.To)
		} else {

			log.Printf("forwarding [%v] to nextHop router #%v\n", rawPacket, nextHop)
			router.addEdgeCostTo(nextHop, packet)
			router.writePacketTo(nextHop, packet)
		}
	}
}

func (router *Router) forwardPacketsFromManager() {
	for {
		data := router.readStringFromManager()
		if data == "QUIT" {
			close(router.managerPacketsDone)
		} else {
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
