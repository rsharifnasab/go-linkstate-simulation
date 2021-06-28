package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Router struct {
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	packetChannel chan string
	Index         int
	Port          int
}

func (router *Router) writeAsString(obj interface{}) {
	router.writer.WriteString(fmt.Sprintf("%v\n", obj))
	router.writer.Flush()
}

func (router *Router) writeAsBytes(obj interface{}) {
	marshalled, err := json.Marshal(obj)
	pnc(err)
	buf := make([]byte, 0)
	buf = append(buf, marshalled...)
	buf = append(buf, '\n')
	_, err = router.writer.Write(buf)
	pnc(err)
	router.writer.Flush()
}

func (router *Router) readString() string {
	str, err := router.reader.ReadString('\n')
	pnc(err)
	return strings.TrimSpace(str)
}

func (router *Router) readInt() int {
	num, err := strconv.Atoi(router.readString())
	pnc(err)
	return num
}

// change name to set connection
func (router *Router) setConnection(conn net.Conn) {
	router.conn = conn
	router.reader = bufio.NewReader(conn)
	router.writer = bufio.NewWriter(conn)
}

func (router *Router) sendDataPackets() {
	for {
		packet := <-router.packetChannel
		log.Printf("router #%v: got packet [%v]\n", router.Index, packet)
		router.writeAsString(packet)
		if packet == "QUIT" {
			break
		}
	}
}
