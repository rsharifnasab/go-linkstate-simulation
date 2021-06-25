package main

import (
	"bufio"
	"fmt"
	"net"
)

var port int
var conn *net.UDPConn
var writer *bufio.Writer

func main() {
	startUDPServer()
	managerConnection, err := net.Dial("tcp", "localhost:8585")
	pnc(err)
	defer managerConnection.Close()
	// reader := bufio.NewReader(managerConnection)
	writer = bufio.NewWriter(managerConnection)
	writeInt(port)
	// for {
	// handleUDPRequest()
	// }
}

func writeInt(num int) {
	writer.WriteString(fmt.Sprintf("%v\n", num))
	defer pnc(writer.Flush())
}

func getSomeFreePort() {
	listener, err := net.Listen("tcp", ":0")
	pnc(err)
	fmt.Printf("Using port: %+v\n", listener.Addr().(*net.TCPAddr))
	port = listener.Addr().(*net.TCPAddr).Port
	pnc(listener.Close())
}

func startUDPServer() (*net.UDPConn, error) {
	var err error
	for failures := 0; failures < 3; failures++ {
		getSomeFreePort()
		addr := net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("127.0.0.1"),
		}
		conn, err := net.ListenUDP("udp", &addr)
		if err == nil {
			return conn, nil
		}
	}
	return nil, err
}

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
