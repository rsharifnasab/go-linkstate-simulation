package main

import (
	"bufio"
	"encoding/json"
	"log"
)

// TODO
// see also sendLSPTo
func (router *Router) recieveLSPs() {
	router.initalCombinedTables()

	log.Printf("(lsp server) ready to get LSPs")
	for i := 0; i < len(router.neighbours); i++ {
		recievedTable := make([]byte, 1000)
		_, _, err := router.conn.ReadFrom(recievedTable[:])
		//log.Printf("conn type is : %T", router.conn)
		reader := bufio.NewReader(router.conn)
		//reader.ReadBytes
		_ = reader
		pnc(err)
		// todo: set recieved table to router.netConns, router.mergedPortMaps
		log.Printf("(lsp server) recieved LSP from router[%v]", -1) // TODO
	}
}

func (router *Router) sendLSPTo(index int) {
	log.Printf("(lsp client) sending LSP to router[%v]\n", index)

	neighboursBytes, err := json.Marshal(router.neighbours)
	pnc(err)
	router.writeUDPAsBytes(index, neighboursBytes)

	portMapBytes, err := json.Marshal(router.portMap)
	pnc(err)
	router.writeUDPAsBytes(index, portMapBytes)

}

func (router *Router) sendLSPs() {
	log.Printf("(lsp client) sending LSP to others {{")
	for _, edge := range router.neighbours {
		router.sendLSPTo(edge.Dest)
	}
	//router.writeToManager("ACKS_RECEIVED")
	log.Printf("}} (lsp client) all lsps send")
}
