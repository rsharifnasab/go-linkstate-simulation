package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	router := &Router{
		doneChannel: make(chan struct{}),
		mpmLock:     sync.RWMutex{},
	}
	defer router.freeResources()

	router.StartUDPServer()
	router.InitLogger()

	router.connectToManager("localhost:8585")
	router.writeToManager(router.port)
	router.getIndexFromManager()
	router.readConnectivityTable()

	router.sendReadySignal()
	router.waitForNetworkSafety()

	router.getPortMap()

	go router.sendAcknowledgements()
	router.testNeighbouringLinks()
	router.waitNetworkReadiness()

	router.initalCombinedTables()
	go router.broadcastSelfLSP()
	router.recieveLSPs()

	router.calculateSPT()

	go router.forwardPacketsFromManager()
	go router.forwardPacketsFromOtherRouters()
	// <-router.doneChannel
	time.Sleep(2 * time.Second)
	log.Printf("router #%v done\n", router.index)
	time.Sleep(time.Second)
}
