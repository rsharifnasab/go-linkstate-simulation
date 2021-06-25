package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
)

func main() {
	log.Println("Server is running")
	manager := createManagerFromConfig("config")
	listener, err := net.Listen("tcp", ":8585")
	pnc(err)
	for i := 0; i < manager.RoutersCount; i++ {
		routerCmd := exec.Command("../router/router")
		reader, err := routerCmd.StderrPipe()
		pnc(err)
		go handleChildError(reader, i)
		routerCmd.Start()

		log.Printf("Started router #%v\n", i)

		conn, err := listener.Accept()
		pnc(err)
		go manager.handleRouter(i, conn)
	}
}

func handleChildError(reader io.ReadCloser, i int) {
	sc := bufio.NewScanner(reader)
	for {
		if !sc.Scan() {
			return
		}
		str := fmt.Sprintf("child %d said: %s\n", i, sc.Text())
		print(str)
	}
}

func pnc(err error) {
	if err != nil {
		panic(err)
	}
}
