package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
)

type LSP struct {
	SenderIndex int
	SenderPort  int

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
	// log.Printf("router[%v] has neighbours : ", index)
	// for _, v := range neighbours {
	// 	log.Printf("  %v", v.Dest)
	// }
	router.netConns[index] = neighbours
}

func (router *Router) receiveSingleLSP() *LSP {
	data := router.readUDPAsBytes()
	lsp := &LSP{}
	pnc(json.Unmarshal(data, lsp))
	return lsp
}

func (router *Router) recieveLSPs() {
	router.initalCombinedTables()

	log.Printf("(lsp server) ready to get LSPs")
	remainingTables := router.routersCount - 1
	isTableReceived := make([]bool, router.routersCount)
	isTableReceived[router.index] = true
	for remainingTables > 0 {
		lsp := router.receiveSingleLSP()
		if !isTableReceived[lsp.SenderIndex] {
			remainingTables--
			log.Printf("(lsp server) recieved LSP from router[%v]", lsp.SenderIndex)
			isTableReceived[lsp.SenderIndex] = true
			router.addToMergedPortMap(lsp.PortMap)
			router.addToNetConns(lsp.SenderIndex, lsp.Neighbours)
			router.broadcastLSP(lsp)
		}
	}
	log.Printf("(lsp server) received the LSPs of all routers")
}

func (router *Router) sendLSPTo(index int, lsp *LSP) {
	log.Printf("(lsp client) sending LSP to router[%v]\n", index)

	bytes, err := json.Marshal(lsp)
	pnc(err)
	router.writeUDPAsBytes(index, bytes)
}

func (router *Router) broadcastLSP(lsp *LSP) {
	log.Printf("(lsp client) sending LSP to others {{")
	for _, edge := range router.neighbours {
		router.sendLSPTo(edge.Dest, lsp)
	}
	//router.writeToManager("ACKS_RECEIVED")
	log.Printf("}} (lsp client) all lsps send")
}

func (router *Router) broadcastSelfLSP() {
	router.broadcastLSP(&LSP{
		SenderIndex: router.index,
		SenderPort:  router.port,

		PortMap:    router.portMap,
		Neighbours: router.neighbours,
	})
}

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func asghar() {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
}
