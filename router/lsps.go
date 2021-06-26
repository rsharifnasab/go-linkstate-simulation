package main

import (
	"encoding/json"
	"log"
)

type LSP struct {
	SenderIndex int
	SenderPort  int

	RecieverIndex int
	RecieverPort  int

	PortMap    map[int]int
	Neighbours []*Edge
}

func (router *Router) addToMergedPortMap(portMap map[int]int) {
	for k, v := range portMap {
		oldVal, isIn := router.mergedPortMaps[k]
		if isIn {
			if oldVal != v {
				log.Fatalf("portmap for[%v] has old val = %v, but i got net val %v", k, oldVal, v)
			}
		}
		router.mergedPortMaps[k] = v
	}
}

func (router *Router) addToNetConns(index int, neighbours []*Edge) {
	log.Printf("router[%v] has neighbours : ", index)
	for _, v := range neighbours {
		log.Printf("  %v", v.Dest)
	}
}

// TODO
// see also sendLSPTo
func (router *Router) recieveLSPs() {
	router.initalCombinedTables()

	log.Printf("(lsp server) ready to get LSPs")
	for i := 0; i < len(router.neighbours); i++ {
		data := router.readUDPAsBytes()
		lsp := &LSP{}
		err := json.Unmarshal(data, lsp)
		pnc(err)
		// todo: set recieved table to router.netConns, router.mergedPortMaps
		router.addToMergedPortMap(lsp.PortMap)
		router.addToNetConns(lsp.SenderIndex, lsp.Neighbours)

		log.Printf("(lsp server) recieved LSP from router[%v]", lsp.SenderIndex) // TODO
	}
}

func (router *Router) sendLSPTo(index int) {
	log.Printf("(lsp client) sending LSP to router[%v]\n", index)

	lsp := &LSP{
		SenderIndex: router.index,
		SenderPort:  router.port,

		RecieverIndex: index,
		RecieverPort:  router.portMap[index],

		PortMap:    router.portMap,
		Neighbours: router.neighbours,
	}
	bytes, err := json.Marshal(lsp)
	pnc(err)
	router.writeUDPAsBytes(index, bytes)
}

func (router *Router) sendLSPs() {
	log.Printf("(lsp client) sending LSP to others {{")
	for _, edge := range router.neighbours {
		router.sendLSPTo(edge.Dest)
	}
	//router.writeToManager("ACKS_RECEIVED")
	log.Printf("}} (lsp client) all lsps send")
}
