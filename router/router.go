package main

import (
	"bufio"
	"net"
	"os"
)

type Router struct {
	conn *net.UDPConn

	port    int
	index   int
	logFile *os.File

	routersCount int
	neighbours   []*Edge
	netConns     [][]Edge

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
