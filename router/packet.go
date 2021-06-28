package main

import (
	"fmt"
	"log"
	"strings"
)

type Packet struct {
	From      int
	To        int
	CostToNow int

	Data string
}

type RawPacket string

func parsePacket(rawPacketNotrim RawPacket) *Packet {
	rawPacket := strings.TrimSpace(string(rawPacketNotrim))
	packet := &Packet{}
	_, err := fmt.Sscanf(rawPacket, "%d %d %d %s",
		&packet.From, &packet.To, &packet.CostToNow, &packet.Data)
	pnc(err)
	return packet
}

func (packet *Packet) serialize() []byte {
	return []byte(fmt.Sprintf("%d %d %d %s",
		packet.From, packet.To, packet.CostToNow, packet.Data))
}

func (router *Router) writePacketTo(nextHop int, packet *Packet) {
	router.writeUDPAsBytes(nextHop, packet.serialize())
}

func (router *Router) updatePacket(nextHop int, packet *Packet) {
	for _, v := range router.neighbours {
		if v.Dest == nextHop {
			packet.CostToNow += v.Cost
			packet.Data = fmt.Sprintf("%s:%d", packet.Data, nextHop)
			return
		}
	}
	log.Fatalf("cannot find next hop cost (%v) for packet (%+v)", nextHop, packet)
}

func (router *Router) sendPacket(rawPacket RawPacket) {

	packet := parsePacket(rawPacket)
	if packet.To == router.index {
		log.Printf("Receive packet from #%v with cost %v. [%v]", packet.From, packet.CostToNow, packet.Data)
	} else {
		nextHop, ok := router.forwardingTable[packet.To]
		if !ok {
			log.Printf("problem forwarding packet (%v) from router #%v to %v\n",
				rawPacket, packet.From, packet.To)
		} else {

			log.Printf("forwarding [%v] to nextHop router #%v\n", rawPacket, nextHop)
			router.updatePacket(nextHop, packet)
			router.writePacketTo(nextHop, packet)
		}
	}
}

// determine packets in udp buffer
// they are from previous broadcasting
func (packet RawPacket) isBad() bool {
	return packet[0] == '{'
}
