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

	portMap        map[int]int
	mergedPortMaps map[int]int
	mpmLock        sync.RWMutex

	managerConnection *net.TCPConn
	managerWriter     *bufio.Writer
	managerReader     *bufio.Reader

	forwardingChannel chan string
	doneChannel       chan struct{}
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
		doneChannel: make(chan struct{}),
		mpmLock:     sync.RWMutex{},
	}
}

func (router *Router) sendPacket(rawPacket string) {
	packet, ok := parsePacket(rawPacket)
	if !ok {
		log.Printf("malformed packet %v at router #%v\n", rawPacket, router.index)
		return
	}
	if packet.destination == router.index {
		log.Printf("Receive packet from #%v. [%v]", packet.source, packet.data)
		return
	}
	nextHop, ok := router.forwardingTable[packet.destination]
	if !ok {
		log.Printf("problem forwarding packet (%v) from router #%v to %v\n", rawPacket, router.index, packet.destination)
		return
	}
	log.Printf("forwarding [%v] to nextHop router #%v\n", rawPacket, nextHop)
	router.writeUDPAsBytes(nextHop, []byte(rawPacket))
}

func (router *Router) forwardPacketsFromManager() {
	for {
		data := router.readStringFromManager()
		if data == "QUIT" {
			close(router.doneChannel)
			return
		}
		//log.Printf("recieved  [%v] from manager\n", data)
		router.sendPacket(data)
	}
}

func (router *Router) forwardPacketsFromOtherRouters() {
	// TODO
	for {
		packet := string(router.readUDPAsBytes())
		_, shouldForward := parsePacket(packet)
		if shouldForward {
			//log.Printf("router #%v received packet: %v\n", router.index, packet)
			router.sendPacket(packet)
			// router.sendPacket(packet)
		} else if packet[0] != '{' {
			log.Printf("router #%v ignored packet: %v\n", router.index, packet)
		}
		router.checkDoneChannel()
	}
}

func (router *Router) checkDoneChannel() {
	select {
	case <-router.doneChannel:
		return
	default:
	}
}
