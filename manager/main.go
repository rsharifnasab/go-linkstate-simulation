package main

import (
	"net"
)

func main() {
	manager := &Manager{}
	manager.loadConfig()
	listener, err := net.Listen("tcp", ":8585")
	pnc(err)
	for i := 0; i < 1; i++ {
		// routerCmd := exec.Command("../router/router")
		// routerCmd.Start()
		// log.Printf("Started router #%v\n", i)
		conn, err := listener.Accept()
		pnc(err)
		manager.handleRouter(i, conn)
	}
}

func pnc(err error) {
	if err != nil {
		panic(err)
	}
}
