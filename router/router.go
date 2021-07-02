package main

import (
	"bufio"
	"net"
	"os"
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

	portMap map[int]int

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

func NewRouter() *Router {
	return &Router{
		managerPacketsDone: make(chan struct{}),
	}
}

func (router *Router) forwardPacketsFromManager() {
	for {
		data := router.readStringFromManager()
		if data == "QUIT" {
			close(router.managerPacketsDone)
		} else {
			rawPacket := RawPacket(data)
			router.sendPacket(rawPacket)
		}
	}
}

func (router *Router) forwardPacketsFromOtherRouters() {
	for {
		packet := RawPacket(router.readUDPAsBytes())
		for packet.isBad() {
			packet = RawPacket(router.readUDPAsBytes())
		}
		router.sendPacket(packet)
	}
}
