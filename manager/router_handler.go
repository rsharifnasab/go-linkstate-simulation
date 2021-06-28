package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// handle our connection to the router in a new go routine
func (manager *Manager) handleRouter(routerIndex int, conn net.Conn) {
	router := manager.routers[routerIndex]
	router.Port = router.readInt()

	log.Printf("router #%v connected, udp port: %v\n",
		router.Index, router.Port)

	router.writeAsString(router.Index)
	router.writeAsString(manager.routersCount)
	router.writeAsBytes(manager.netConns[router.Index])
	log.Printf("connectivity table sent for router[%v]", router.Index)

	manager.getReadySignalFromRouter(router) // wait for this router
	<-manager.readyChannel                   // wait for other routers
	router.writeAsString("safe")             // tell routers they can ack
	manager.sendPortMap(router)              // send all routers ports to router

	manager.getAcksReceivedFromRouter(router) // my router got portmap
	<-manager.networkReadyChannel             // all routers got portmap

	router.writeAsString("NETWORK_READY")
	router.handlePackets()
}

func (manager *Manager) getReadySignalFromRouter(router *Router) {
	readiness := router.readString()
	if readiness == "READY" {
		manager.readyWG.Done()
	} else {
		panic("Router couldn't get ready.")
	}
}

func (manager *Manager) sendPortMap(router *Router) {
	portMap := make(map[int]int)
	for _, edge := range manager.netConns[router.Index] {
		portMap[edge.Dest] = manager.routers[edge.Dest].Port
	}
	router.writeAsBytes(portMap)
}

func (manager *Manager) getAcksReceivedFromRouter(router *Router) {
	str := router.readString()
	if str != "ACKS_RECEIVED" {
		panic(fmt.Sprintf("router #%v didn't receive acks: %v", router.Index, str))
	}
	manager.networkReadyWG.Done()
}

func handleChildError(reader io.ReadCloser, i int) {
	sc := bufio.NewScanner(reader)
	for {
		if !sc.Scan() {
			return
		}
		log.Printf("child %d : %s\n", i, sc.Text())
	}
}
