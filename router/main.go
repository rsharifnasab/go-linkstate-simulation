package main

import "time"

func main() {
	router := &Router{}
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

	go router.recieveLSPs()
	router.sendLSP()

	time.Sleep(10 * time.Second)
}
