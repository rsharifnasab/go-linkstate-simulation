package main

import (
	"encoding/json"
	"log"
)

type LSP struct {
	SenderIndex int
	SenderPort  int

	PortMap    map[int]int
	Neighbours []*Edge
}

func (router *Router) addToMergedPortMap(portMap map[int]int) {
	router.mpmLock.Lock()
	defer router.mpmLock.Unlock()
	for k, v := range portMap {
		oldVal, isIn := router.mergedPortMaps[k]
		if isIn {
			if oldVal != v {
				log.Fatalf("portmap for[%v] has old val = %v, but i got net val %v", k, oldVal, v)
			}
		}
		if v != -1 {
			router.mergedPortMaps[k] = v
		}
	}
}

func (router *Router) addToNetConns(index int, neighbours []*Edge) {
	// log.Printf("router[%v] has neighbours : ", index)
	// for _, v := range neighbours {
	// 	log.Printf("  %v", v.Dest)
	// }
	router.netConns[index] = neighbours
}

func (router *Router) receiveSingleLSP() *LSP {
	data := router.readUDPAsBytes()
	lsp := &LSP{}
	// log.Printf("received %v as lsp\n", strings.TrimSpace(string(data)))
	// pnc(json.Unmarshal([]byte(strings.TrimSpace(string(data))), lsp))
	pnc(json.Unmarshal(data, lsp))
	return lsp
}

func (router *Router) recieveLSPs() {

	log.Printf("(lsp server) listening to other routers LSPs started")
	remainingTables := router.routersCount - 1
	isTableReceived := make([]bool, router.routersCount)
	isTableReceived[router.index] = true
	for remainingTables > 0 {
		lsp := router.receiveSingleLSP()
		if !isTableReceived[lsp.SenderIndex] {
			remainingTables--
			log.Printf("(lsp server) recieved #%v LSP", lsp.SenderIndex)
			isTableReceived[lsp.SenderIndex] = true
			router.addToMergedPortMap(lsp.PortMap)
			router.addToNetConns(lsp.SenderIndex, lsp.Neighbours)
			router.broadcastLSP(lsp)
		}
	}
	log.Printf("(lsp server) receiving LSPs done")
}

func (router *Router) sendLSPTo(index int, lsp *LSP) {
	log.Printf("(lsp client) sending #%v LSP to router[%v]\n", lsp.SenderIndex, index)

	bytes, err := json.Marshal(lsp)
	pnc(err)
	router.writeUDPAsBytes(index, bytes)
}

func (router *Router) broadcastLSP(lsp *LSP) {
	log.Printf("(lsp client) broadcasting #%v LSP to neighbours", lsp.SenderIndex)
	for _, edge := range router.neighbours {
		router.sendLSPTo(edge.Dest, lsp)
	}
	//log.Printf("(lsp client)  LSPs done")
}

func (router *Router) broadcastSelfLSP() {
	router.broadcastLSP(&LSP{
		SenderIndex: router.index,
		SenderPort:  router.port,

		PortMap:    router.portMap,
		Neighbours: router.neighbours,
	})
}
