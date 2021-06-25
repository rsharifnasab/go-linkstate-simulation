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
	manager := createManagerFromConfig("config")
	listener, err := net.Listen("tcp", ":8585")
	pnc(err)
	for i := 0; i < manager.RoutersCount; i++ {
		routerCmd := exec.Command("../router/router")
		reader, err := routerCmd.StderrPipe()
		pnc(err)
		go handleChildError(reader, i)
		routerCmd.Start()

		log.Printf("Created router #%v\n", i)

		conn, err := listener.Accept()
		pnc(err)
		go manager.handleRouter(i, conn)
	}
	time.Sleep(1 * time.Second)
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
