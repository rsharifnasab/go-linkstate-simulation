package main

import (
	"log"
	"net"
	"os/exec"
)

func main() {
	initLogger()
	// go sniff()
	manager := newManagerWithConfig("config")
	listener, err := net.Listen("tcp", ":8585")
	pnc(err)
	for i := 0; i < manager.routersCount; i++ {
		routerCmd := exec.Command("../router/router")
		reader, err := routerCmd.StderrPipe()
		pnc(err)
		go handleChildError(reader, i)
		routerCmd.Start()

		log.Printf("router #%v: created\n", i)
		conn, err := listener.Accept()
		pnc(err)
		manager.routers[i].setConnection(conn)
		manager.readyWG.Add(1)
		go manager.handleRouter(i, conn)
	}
	log.Printf("waiting for routers to get ready")
	manager.readyWG.Wait()
	close(manager.readyChannel)
	log.Printf("all routers ready. routers checking direct links")
	manager.networkReadyWG.Add(manager.routersCount)
	manager.networkReadyWG.Wait()
	close(manager.networkReadyChannel)
	log.Printf("Network is ready")
	manager.sendTestPackets()
	for _, router := range manager.routers {
		router.conn.Close()
	}
}

// func StartUDPServer() {
// 	const MAX_TRIES = 3
// 	var err error
// 	port := "6868"
// 	// log.Printf("getSomeFreePort() provided port number %v\n", port)
// 	addr := fmt.Sprintf(":%d", port)
// 	conn, err := net.ListenPacket("udp", addr)
// 	if err == nil {
// 		conn = conn.(*net.UDPConn)
// 		connReader = bufio.NewReader(router.conn)
// 		connWriter = bufio.NewWriter(router.conn)
// 		port = port
// 	}
// 	// router.connWriter.Write([]byte("salam"))
// }
