package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

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

func (router *Router) readStringFromManager() string {
	str, err := router.managerReader.ReadString('\n')
	if errors.Is(err, io.EOF) {
		return "QUIT"
	}
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
	log.Printf("reading connectivity table")

	router.routersCount = router.readIntFromManager()
	rawMessage := router.readBytesFromManager()
	pnc(json.Unmarshal(rawMessage, &router.neighbours))
	for _, edge := range router.neighbours {
		log.Printf("  {Dest: %+v, Cost: %v}\n", edge.Dest, edge.Cost)
	}
}

func (router *Router) sendReadySignal() {
	router.writeToManager("READY")
	log.Printf("I am ready")
}

func (router *Router) waitForNetworkSafety() {
	message := router.readStringFromManager()
	if message != "safe" {
		log.Fatal("We are not safe")
	}
	log.Printf("We are all safe")
}

func (router *Router) getPortMap() {
	rawMessage := router.readBytesFromManager()
	pnc(json.Unmarshal(rawMessage, &router.portMap))
}
