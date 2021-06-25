package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Router struct {
	conn    *net.UDPConn
	port    int
	logFile *os.File

	routersCount      int
	neighbours        []*Edge
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
		addr := net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("127.0.0.1"),
		}
		conn, err := net.ListenUDP("udp", &addr)
		if err == nil {
			router.conn = conn
			router.port = port
			return
		}
	}
	pnc(err)
}

func (router *Router) readConnTable() {
	log.Printf("read conn table")
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
		panic(err)
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

func (router *Router) readIntFromManager() int {
	str, err := router.managerReader.ReadString('\n')
	pnc(err)
	num, err := strconv.Atoi(strings.TrimSpace(str))
	pnc(err)
	return num
}

func (router *Router) readConnectivityTable() {
	router.routersCount = router.readIntFromManager()
	log.Printf("received %v. waiting for connectivity table", router.routersCount)
	rawMessage, err := router.managerReader.ReadBytes('\n')
	pnc(err)
	pnc(json.Unmarshal(rawMessage, &router.neighbours))
	for _, edge := range router.neighbours {
		log.Printf("%+v\n", edge)
	}
}

func (router *Router) sendReadySignal() {
	// for debug
	//time.Sleep(5 * time.Second)

	router.writeToManager("READY")
	log.Printf("I am ready")
}

func (router *Router) waitForOurRouters() {
	router.readIntFromManager()
	log.Printf("we are all synced")
}
