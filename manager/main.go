package main

import (
	"log"
	"time"
)

func main() {
	initLogger()

	manager := newManagerWithConfig("config")
	defer manager.freeResources()

	for i := 0; i < manager.routersCount; i++ {
		manager.launchRouter(i)
	}

	log.Printf("Manager: Waiting for routers to get ready")
	manager.readyWG.Wait()
	close(manager.readyChannel)
	log.Printf("Manager: All routers ready. routers checking direct links")

	manager.networkReadyWG.Wait()
	close(manager.networkReadyChannel)
	log.Printf("Manager: Network is ready")

	// send tests in yaml to source routers
	// they will do the rest of the job
	manager.sendTestPackets()

	// wait for routers to transmit packets
	// if we quit now, they will be terminated by os
	time.Sleep(5 * time.Second)

}
