package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func pnc(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	port := 4567
	log.Printf("dialing")
	conn, err := net.Dial("udp", fmt.Sprintf("localhost:%v", port))
	pnc(err)
	log.Printf("dial complete")

	writer := bufio.NewWriter(conn)
	writer.WriteString(fmt.Sprintf("%v\n", "my index"))
	writer.Flush()
	log.Printf("sent ack request to %v on %v\n", "my index", port)
}
