package main

import (
	"log"
	"time"
)

func main() {
	router := NewRouter()
	defer router.freeResources()

	router.StartUDPServer()

	router.connectToManager("localhost:8585")
	router.writeToManager(router.port)
	router.getIndexFromManager()

	router.InitLogger()

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
	time.Sleep(3 * time.Second)
	log.Printf("done\n")
}
