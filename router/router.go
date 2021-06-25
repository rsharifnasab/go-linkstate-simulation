package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var ManagerConn *net.UDPConn
var managerWriter *bufio.Writer
var managerReader *bufio.Reader

func main() {
	udpConn, err, port := startUDPServer()
	pnc(err)
	logFile := initloger(port)
	defer logFile.Close()

	managerConnection, err := net.Dial("tcp", "localhost:8585")
	pnc(err)
	defer managerConnection.Close()

	managerReader = bufio.NewReader(managerConnection)
	managerWriter = bufio.NewWriter(managerConnection)
	managerWrite(port)

	// todo
	readConnTable()

	_ = udpConn
}

func readConnTable() {
	log.Printf("read conn table")
}

func initloger(port int) *os.File {
	logFileAdd := fmt.Sprintf("../%v.log", port)
	logFile, err := os.OpenFile(logFileAdd, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	pnc(err)

	log.SetOutput(logFile)
	log.SetFlags(0)
	log.Printf("")
	//log.Printf("- - - %v logger - - -", port)
	log.SetPrefix(fmt.Sprintf("child %v ", port))

	return logFile
}

func managerWrite(data interface{}) {
	managerWriter.WriteString(fmt.Sprintf("%v\n", data))
	pnc(managerWriter.Flush())
}

func getSomeFreePort() int {
	listener, err := net.Listen("tcp", ":0")
	pnc(err)
	//fmt.Fprintf(os.Stderr, "using port: %+v\n", listener.Addr().(*net.TCPAddr))
	pnc(listener.Close())
	return listener.Addr().(*net.TCPAddr).Port
}

func startUDPServer() (*net.UDPConn, error, int) {
	var err error
	for failures := 0; failures < 3; failures++ {
		port := getSomeFreePort()
		addr := net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("127.0.0.1"),
		}
		conn, err := net.ListenUDP("udp", &addr)
		if err == nil {
			return conn, nil, port
		}
	}
	return nil, err, -1
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
