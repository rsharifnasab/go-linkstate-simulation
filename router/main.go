package main

import "time"

func main() {
	router := &Router{}
	defer router.freeResources()
	router.StartUDPServer()
	router.InitLogger()
	router.connectToManager("localhost:8585")
	router.writeToManager(router.Port)
	router.getIndexFromManager()
	router.readConnectivityTable()
	router.sendReadySignal()
	router.waitForNetworkSafety()
	// router.testNeighbouringLinks()
	// router.waitNetworkReadiness()
	time.Sleep(10 * time.Second)
}
