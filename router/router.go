package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

type Router struct {
	conn net.PacketConn

	port    int
	index   int
	logFile *os.File

	routersCount      int
	neighbours        []*Edge
	portMap           map[int]int
	managerConnection *net.TCPConn
	managerWriter     *bufio.Writer
	managerReader     *bufio.Reader
}

type Edge struct {
	Dest int
	Cost int
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
	log.SetPrefix(fmt.Sprintf("child %v ", router.port))

}

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
			router.conn = conn
			router.port = port
			break
		}
	}
	pnc(err)
}

func getSomeFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	pnc(err)
	//fmt.Fprintf(os.Stderr, "using port: %+v\n", listener.Addr().(*net.TCPAddr))
	pnc(listener.Close())
	return listener.Addr().(*net.TCPAddr).Port
}

func pnc(err error) {
	if err != nil {
		log.Writer().Write(debug.Stack())
		log.Fatal(err)
	}
}
func (router *Router) connectToManager(add string) {
	connection, err := net.Dial("tcp", "localhost:8585")
	pnc(err)
	router.managerConnection = connection.(*net.TCPConn)
	router.managerReader = bufio.NewReader(router.managerConnection)
	router.managerWriter = bufio.NewWriter(router.managerConnection)
}

func (router *Router) writeToManager(data interface{}) {
	router.managerWriter.WriteString(fmt.Sprintf("%v\n", data))
	pnc(router.managerWriter.Flush())
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

func (router *Router) readStringFromManager() string {
	str, err := router.managerReader.ReadString('\n')
	pnc(err)
	return strings.TrimSpace(str)
}

func (router *Router) readBytesFromManager() []byte {
	rawMessage, err := router.managerReader.ReadBytes('\n')
	pnc(err)
	return rawMessage[:len(rawMessage)-1]
}

func (router *Router) readIntFromManager() int {
	num, err := strconv.Atoi(router.readStringFromManager())
	pnc(err)
	return num
}

func (router *Router) getIndexFromManager() {
	router.index = router.readIntFromManager()
}

func (router *Router) readConnectivityTable() {
	router.routersCount = router.readIntFromManager()
	log.Printf("router #%v waking up", router.index)
	rawMessage, err := router.managerReader.ReadBytes('\n')
	pnc(err)
	pnc(json.Unmarshal(rawMessage, &router.neighbours))
	for _, edge := range router.neighbours {
		log.Printf("{Dest: %+v, Cost: %v}\n", edge.Dest, edge.Cost)
	}
}

func (router *Router) sendReadySignal() {
	// for debug
	//time.Sleep(5 * time.Second)

	router.writeToManager("READY")
	log.Printf("I am ready")
}

func (router *Router) waitForNetworkSafety() {
	message := router.readStringFromManager()
	if message != "safe" {
		panic("we are not safe")
	}
	log.Printf("we are all safe")
}

func (router *Router) getPortMap() {
	rawMessage := router.readBytesFromManager()
	pnc(json.Unmarshal(rawMessage, &router.portMap))
}

func (router *Router) testNeighbouringLinks() {
	log.Printf("checking neighbouring links")
	for _, edge := range router.neighbours {
		index := edge.Dest
		port := router.portMap[index]
		log.Printf("dialing to router[%v] on port %v\n", index, port)
		conn := dialUDP(fmt.Sprintf("localhost:%v", port))
		router.sendAckRequest(conn, index, port)
		router.getAckResponse(conn, index, port)
		conn.Close()
		//log.Printf("%v check", edge.Dest)
	}
	router.writeToManager("ACKS_RECEIVED")
}

func (router *Router) sendAckRequest(conn net.Conn, index, port int) {
	writer := bufio.NewWriter(conn)
	writer.WriteString(fmt.Sprintf("%v\n", router.index))
	writer.Flush()
	log.Printf("sent ack request to %v on %v\n", index, port)
}

func (router *Router) getAckResponse(conn net.Conn, index, port int) {
	ackResponse, err := bufio.NewReader(conn).ReadString('\n')
	pnc(err)
	if ackResponse != "ack\n" {
		panic(fmt.Sprintf("Who are you not to acknowledge me router #%v listening on port %v by saying %v", index, port, ackResponse))
	}
	log.Printf("received ack from %v on %v\n", index, port)

}

func dialUDP(addr string) net.Conn {
	conn, err := net.Dial("udp", addr)
	pnc(err)
	return conn
}

func (router *Router) sendAcknowledgements() {
	log.Printf("(udp server) listening for other routers")
	for i := 0; i < len(router.neighbours); i++ {
		ackRequest := make([]byte, 100)
		n, addr, err := router.conn.ReadFrom(ackRequest[:])
		pnc(err)
		//log.Printf("(udp server) ack req from  router[%v]", string(ackRequest[:n-1]))
		router.conn.WriteTo([]byte("ack\n"), addr)
		log.Printf("(udp server) acknowledged  router[%v]", string(ackRequest[:n-1]))
	}
}
