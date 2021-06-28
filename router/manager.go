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
	rawMessage := router.readBytesFromManager()
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
		log.Fatal("we are not safe")
	}
	log.Printf("we are all safe")
}

func (router *Router) getPortMap() {
	rawMessage := router.readBytesFromManager()
	pnc(json.Unmarshal(rawMessage, &router.portMap))
}