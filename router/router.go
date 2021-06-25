package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Router struct {
	conn    *net.UDPConn
	Port    int
	logFile *os.File

	RouterCount int
	Neighbors   []*Edge
}

type Edge struct {
	Dist int
	Cost int
}

func (router *Router) InitLogger() {
	logFileAdd := fmt.Sprintf("../%v.log", router.Port)
	logFile, err := os.OpenFile(logFileAdd, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	pnc(err)

	router.logFile = logFile

	log.SetOutput(logFile)
	log.SetFlags(0)
	log.Printf("")
	//log.Printf("- - - %v logger - - -", port)
	log.SetPrefix(fmt.Sprintf("child %v ", router.Port))

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
			router.Port = port
		}
	}
	panic(err)
}

func main() {
	router := &Router{}

	router.StartUDPServer()

	router.InitLogger()
	defer router.logFile.Close()

	router.ConnectToManager("localhost:8585")
	defer router.manager.conn.Close()

	manager.write(port)

	// todo
	newRouter(manger)

}

func readConnTable() {
	log.Printf("read conn table")
}

func getSomeFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	pnc(err)
	//fmt.Fprintf(os.Stderr, "using port: %+v\n", listener.Addr().(*net.TCPAddr))
	pnc(listener.Close())
	return listener.Addr().(*net.TCPAddr).Port
}

// unused by now
func udpClient() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}

func pnc(err error) {
	if err != nil {
		panic(err)
	}
}
