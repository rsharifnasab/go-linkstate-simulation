package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"
)

func pnc(err error) {
	if err != nil {
		log.Writer().Write(debug.Stack())
		log.Fatal(err)
	}
}

func getSomeFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	pnc(err)
	//fmt.Fprintf(os.Stderr, "using port: %+v\n", listener.Addr().(*net.TCPAddr))
	pnc(listener.Close())
	return listener.Addr().(*net.TCPAddr).Port
}

func (router *Router) InitLogger() {
	logFileAdd := fmt.Sprintf("../%v.log", router.port)
	logFile, err := os.OpenFile(logFileAdd, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	pnc(err)

	router.logFile = logFile

	log.SetOutput(logFile)
	log.SetFlags(0)
	log.Printf("")
	//log.Printf("- - - %v logger - - -", port)
	log.SetPrefix(fmt.Sprintf("router on %v ", router.port))

}

func (router *Router) freeResources() {
	if router.logFile != nil {
		router.logFile.Close()
	}
	if router.managerConnection != nil {
		router.managerConnection.Close()
	}
	if router.conn != nil {
		router.conn.Close()
	}
}

// open udp port as client
func dialUDP(addr string) net.Conn {
	conn, err := net.Dial("udp", addr)
	pnc(err)
	return conn
}

func (router *Router) readUDPAsBytes() []byte {
	buff := make([]byte, 1024*1024)
	n, _, err := router.conn.ReadFrom(buff)
	pnc(err)

	return buff[:n]
}

// open an udp socket and send byte slice to specified router
func (router *Router) writeUDPAsBytes(index int, data []byte) {
	port := router.portMap[index]
	conn := dialUDP(fmt.Sprintf("localhost:%v", port))
	defer conn.Close()

	//toSendData := make([]byte, 0)
	//toSendData = append(toSendData, data...)
	//toSendData = append(toSendData, '\n')
	//_, err := conn.Write(toSendData)

	_, err := conn.Write(data)
	pnc(err)
}

// open a UDP server on desired port and start listening
func (router *Router) StartUDPServer() {
	const MAX_TRIES = 3
	var err error
	for failures := 0; failures < MAX_TRIES; failures++ {
		port := getSomeFreePort()
		// log.Printf("getSomeFreePort() provided port number %v\n", port)
		if port == 0 {
			log.Fatal("system provided port number 0")
			continue
		}
		addr := fmt.Sprintf(":%d", port)
		conn, err := net.ListenPacket("udp", addr)
		if err == nil {
			router.conn = conn.(*net.UDPConn)
			router.connReader = bufio.NewReader(router.conn)
			router.connWriter = bufio.NewWriter(router.conn)
			router.port = port
			break
		}
		router.connWriter.Write([]byte("salam"))
	}
	pnc(err)
}

// we need another table near or primary table
// for keep all connections in network, not just ours.
// this function will initialize it from primary table
func (router *Router) initalCombinedTables() {
	router.netConns = make([][]*Edge, router.routersCount)
	router.netConns[router.index] = router.neighbours
	router.mergedPortMaps = make(map[int]int)
	for k, v := range router.portMap {
		router.mergedPortMaps[k] = v
	}
}

func createSlice(size int, defaultValue int) []int {
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = defaultValue
	}
	return slice
}
