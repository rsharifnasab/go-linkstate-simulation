package main

import (
	"container/heap"
	"log"
	"math"
)

func (router *Router) calculateSPT() {
	pq := make(PriorityQueue, 0)
	mark := make([]bool, router.routersCount)
	prev := createSlice(router.routersCount, -1)
	shortestPathCost := createSlice(router.routersCount, math.MaxInt32)
	shortestPathCost[router.index] = 0
	router.initPQ(&pq)
	router.forwardingTable = make(map[int]int)
	router.forwardingTable[router.index] = -1
	for pq.Len() > 0 {
		pqItem := heap.Pop(&pq).(*PriorityQueueItem)
		for mark[pqItem.index] && pq.Len() > 0 {
			pqItem = heap.Pop(&pq).(*PriorityQueueItem)
			if !mark[pqItem.index] {
				break
			}
		}
		if mark[pqItem.index] {
			break
		}
		mark[pqItem.index] = true
		for _, edge := range router.netConns[pqItem.index] {
			if mark[edge.Dest] {
				continue
			}
			newCost := pqItem.dist + edge.Cost
			if newCost >= shortestPathCost[edge.Dest] {
				continue
			}
			shortestPathCost[edge.Dest] = newCost
			// push a new item to queue with path source -> pqItem.index -> edge.Dest
			heap.Push(&pq, &PriorityQueueItem{
				dist:  newCost,
				index: edge.Dest,
			})
			// update forwarding table
			firstRouterInPath := router.forwardingTable[pqItem.index]
			if firstRouterInPath == -1 {
				firstRouterInPath = edge.Dest
			}
			router.forwardingTable[edge.Dest] = firstRouterInPath
			// update shortest path tree
			prev[edge.Dest] = pqItem.index
		}
	}
	log.Printf("Calculated shortest Path Tree:\n\t%+v", router.forwardingTable)

}

func (router *Router) initPQ(pq *PriorityQueue) {
	heap.Init(pq)
	isNeighbour := make([]bool, router.routersCount)
	heap.Push(pq, &PriorityQueueItem{
		dist:  0,
		index: router.index,
	})
	for _, edge := range router.netConns[router.index] {
		heap.Push(pq, &PriorityQueueItem{
			dist:  edge.Cost,
			index: edge.Dest,
		})
		isNeighbour[edge.Dest] = true
	}
}
