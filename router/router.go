package main

import (
	"bufio"
	"net"
	"os"
	"time"
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

	managerConnection *net.TCPConn
	managerWriter     *bufio.Writer
	managerReader     *bufio.Reader
}

type Edge struct {
	Dest int
	Cost int
}

func (router *Router) sendPacket(packet string) {

}

func (router *Router) sendPacketsGotFromManager() {
	for {
		data := router.readStringFromManager()
		if data == "QUIT" {
			time.Sleep(1 * time.Second)
			os.Exit(0)
		}

		router.sendPacket(data)

	}
}
