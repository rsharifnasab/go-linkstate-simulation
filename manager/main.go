package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
	"time"
)

func main() {
	initLogger()
	manager := newManagerWithConfig("config")
	listener, err := net.Listen("tcp", ":8585")
	pnc(err)
	conns := make([]*net.Conn, 0)
	for i := 0; i < manager.routersCount; i++ {
		routerCmd := exec.Command("../router/router")
		reader, err := routerCmd.StderrPipe()
		pnc(err)
		go handleChildError(reader, i)
		routerCmd.Start()

		log.Printf("Created router #%v\n", i)

		conn, err := listener.Accept()
		pnc(err)
		conns = append(conns, &conn)
		manager.readyWG.Add(1)
		go manager.handleRouter(i, conn)
	}
	log.Printf("waiting for routers to get ready")
	manager.readyWG.Wait()
	close(manager.readyChannel)
	log.Printf("all routers got ready")

	time.Sleep(3 * time.Second)
	for _, conn := range conns {
		(*conn).Close()
	}
}

func initLogger() {
	log.SetFlags(0)
	log.Println("Server is running")
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

func pnc(err error) {
	if err != nil {
		panic(err)
	}
}
