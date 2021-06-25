package main

func main() {
	router := &Router{}
	defer router.freeResources()
	router.StartUDPServer()
	router.InitLogger()
	router.connectToManager("localhost:8585")
	router.writeToManager(router.port)
	router.readConnectivityTable()
	router.sendReadySignal()
	router.waitForOurRouters()
}
