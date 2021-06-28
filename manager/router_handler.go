package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
)

// handle our connection to the router in a new go routine
func (manager *Manager) handleRouter(routerIndex int, conn net.Conn) {
	router := manager.routers[routerIndex]
	router.Port = router.readInt()

	log.Printf("router #%v: connected, UDP port: %v\n",
		router.Index, router.Port)

	router.writeAsString(router.Index)
	router.writeAsString(manager.routersCount)
	router.writeAsBytes(manager.netConns[router.Index])
	log.Printf("router #%v: connectivity table sent", router.Index)

	manager.getReadySignalFromRouter(router) // wait for this router
	<-manager.readyChannel                   // wait for other routers
	router.writeAsString("safe")             // tell routers they can ack
	manager.sendPortMap(router)              // send all routers ports to router

	manager.getAcksReceivedFromRouter(router) // my router got portmap
	<-manager.networkReadyChannel             // all routers got portmap

	router.writeAsString("NETWORK_READY")

	router.sendDataPackets()
}

func (router *Router) sendDataPackets() {
	for {
		// read from packet queue
		packet := <-router.packetChannel
		log.Printf("router #%v: got packet [%v]\n", router.Index, packet)
		router.writeAsString(packet)

		if packet == "QUIT" {
			// our last packet
			break
		}
	}
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
		log.Printf("router #%v (STDERR) : %s\n", i, sc.Text())
	}
}

func (manager *Manager) launchRouter(i int) {
	routerCmd := exec.Command("../router/router")
	reader, err := routerCmd.StderrPipe()
	pnc(err)

	go handleChildError(reader, i)
	routerCmd.Start()

	log.Printf("router #%v: created\n", i)
	conn, err := manager.listener.Accept()
	pnc(err)
	manager.routers[i].setConnection(conn)
	go manager.handleRouter(i, conn)
}
